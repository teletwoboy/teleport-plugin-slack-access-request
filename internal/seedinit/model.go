package seedinit

import "time"

type SeedInit struct {
	SeedInitId int32
	Status     string
	UseYn      bool
	CreateCode string
	CreateDate time.Time
	UpdateCode string
	UpdateDate time.Time
	DeleteCode string
	DeleteDate time.Time
	Version    int64
}
