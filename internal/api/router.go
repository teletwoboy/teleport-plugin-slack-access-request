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

package api

import (
	"net/http"
	v1 "teleport-plugin-slack-access-request/internal/api/v1"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	v1 *v1.Router
}

func NewRouter(v1 *v1.Router) *Router {
	return &Router{
		v1: v1,
	}
}

func (r *Router) Setup(router *chi.Mux) http.Handler {
	router.With(VerifySlackRequest()).
		Route("/api", func(router chi.Router) {
			r.v1.Route(router)
		})
	return router
}
