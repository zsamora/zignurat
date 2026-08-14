package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Names     string
	Surnames  string
	DateBirth time.Time
	UUID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();<-:create"`
}

type Organization struct {
	gorm.Model
	OrgType uint8
	Name    string
	Diocese string
	Commune string
	Address string
	AdminID uint
	UUID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();<-:create"`
}
