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

package database

import (
	"time"
)

const (
	CreateCode = "teleport-plugin-slack-access-request"
	UpdateCode = "teleport-plugin-slack-access-request"
	DeleteCode = "teleport-plugin-slack-access-request"
)

type BaseEntity struct {
	UseYn      bool
	CreateCode string
	CreateDate time.Time
	UpdateCode string
	UpdateDate time.Time
	DeleteCode string
	DeleteDate time.Time
	Version    int64
}

func MarkCreate() *BaseEntity {
	return &BaseEntity{
		UseYn:      true,
		CreateCode: CreateCode,
		CreateDate: time.Now().UTC().Truncate(time.Second),
		Version:    0,
	}
}

func MarkUpdate() *BaseEntity {
	return &BaseEntity{
		UpdateCode: UpdateCode,
		UpdateDate: time.Now().UTC().Truncate(time.Second),
	}
}

func MarkDelete() *BaseEntity {
	return &BaseEntity{
		UseYn:      false,
		DeleteCode: DeleteCode,
		DeleteDate: time.Now().UTC().Truncate(time.Second),
	}
}
