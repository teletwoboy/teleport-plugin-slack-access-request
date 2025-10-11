/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"teleport-plugin-slack-access-request/internal/api/check"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/metric"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/worker"
	"teleport-plugin-slack-access-request/internal/util/container"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	router := chi.NewRouter()
	isReady := &atomic.Value{}
	isReady.Store(false)

	app := NewContext()
	setupCloseHandler(cancel, func() {
		app.Cleanup(ctx)
	})

	errCh := make(chan error, 1)
	go func() {
		if err := startCheckServer(router, isReady, app); err != nil {
			errCh <- err
		}
	}()
	go func() {
		if err := startAPIServer(ctx, router, isReady, app); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		slog.Error("Error starting API server", "err", err)
		app.Cleanup(ctx)
		cancel()
	}
}

func setupCloseHandler(cancel context.CancelFunc, cleanup func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		slog.Info("Caught sig interrupt...exiting.")
		cleanup()
		cancel()
	}()
}

func startCheckServer(router *chi.Mux, isReady *atomic.Value, app *Context) error {
	router.Use(middleware.Recoverer)

	router.Get("/healthz", check.Healthz)
	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		check.Readyz(w, r, isReady)
	})

	srv := &http.Server{
		Addr:    ":" + config.Cfg.Server.Port,
		Handler: router,
	}
	app.Server = srv

	slog.Info("starting http server", "port", config.Cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	return nil
}

func startAPIServer(ctx context.Context, router *chi.Mux, isReady *atomic.Value, app *Context) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	slog.Info("successfully connected to database")
	app.DB = db

	clients, err := container.NewClients(ctx)
	if err != nil {
		return fmt.Errorf("failed to create clients: %w", err)
	}
	slog.Info("successfully initialized clients")
	app.Clients = clients

	repos := container.NewRepositories(db.Queries)
	services := container.NewServices(clients, repos)

	if err := services.SeedInit.Init(ctx, db, clients.Slack, clients.Teleport); err != nil {
		return fmt.Errorf("failed to seed init: %w", err)
	}

	metric.Init(db)
	slog.Info("successfully initialized metric for prometheus")

	if config.Cfg.Otel.Enable {
		tpShutdown, err := telemetry.Init(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize telemetry: %w", err)
		}
		slog.Info("successfully initialized open telemetry")
		app.OtelShutdown = tpShutdown
	}

	event := NewEvent(db, clients, services)
	go event.StartWatcher(ctx)
	slog.Info("starting event watching")

	go worker.StartWorker(ctx, db, clients, services)
	slog.Info("starting outbox worker")

	routers := NewRouter(db, clients, repos, services)
	routers.Setup(router)
	router.Handle("/metrics", promhttp.HandlerFor(metric.Registry, promhttp.HandlerOpts{}))
	isReady.Store(true)
	return nil
}
