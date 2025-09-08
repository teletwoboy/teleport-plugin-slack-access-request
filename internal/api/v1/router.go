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

package v1

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	AccessPolicyHandler *OpenAccessPolicyModalHandler
	AccessRoleHandler   *OpenAccessRequestModalHandler
	InteractionHandler  *InteractionHandler
}

func NewRouter(ap *OpenAccessPolicyModalHandler, ar *OpenAccessRequestModalHandler, i *InteractionHandler) *Router {
	return &Router{
		AccessPolicyHandler: ap,
		AccessRoleHandler:   ar,
		InteractionHandler:  i,
	}
}

func (r *Router) Route(router chi.Router) {
	router.Route("/v1", func(router chi.Router) {
		router.Post("/access-policy", r.AccessPolicyHandler.Handle)
		router.Post("/access-request", r.AccessRoleHandler.Handle)
		router.Post("/interaction", r.InteractionHandler.Handle)
	})
}
