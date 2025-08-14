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

package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func Init() {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get current working directory: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.SourceKey:
				src := a.Value.Any().(*slog.Source)
				rel := strings.TrimPrefix(src.File, rootDir)
				rel = strings.TrimPrefix(rel, string(os.PathSeparator))
				return slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", rel, src.Line))
			case slog.TimeKey:
				t := a.Value.Time()
				return slog.String(slog.TimeKey, t.Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}))
	slog.SetDefault(logger) // set Global
}
