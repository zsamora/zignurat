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
	Books   []Book    `gorm:"foreignKey:OrgID"`
	UUID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();<-:create"`
}

type Book struct {
	gorm.Model
	OrgID            uint
	BookNr           uint8
	DateInitial      time.Time
	DateFinal        time.Time
	RegistersBaptism []RegisterBaptism `gorm:"foreignKey:BookID"`
	IndexBaptism     []IndexBaptism    `gorm:"foreignKey:BookID"`
	UUID             uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();<-:create"`
	BookType         uint8
}

type RegisterBaptism struct {
	gorm.Model
	BookID        uint
	PageNumber    uint16
	RegNumber     uint16
	IndexID       uint
	OrgBaptism    uint
	Baptizer      string
	DateBaptism   time.Time
	BaptizedNameF string
	BaptizedNameS string
	RUT           string
	DateBirth     time.Time
	PlaceBirth    string
	FatherName    string
	FatherSurname string
	MotherName    string
	MotherSurname string
	Godfather     string
	Godmother     string
	UUID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();<-:create"`
}

type IndexBaptism struct {
	gorm.Model
	OrgID          uint
	BookID         uint
	RegID          uint
	UserSurnameF   string
	UserSurnameM   string
	UserNameFirst  string
	UserNameSecond string
	PageNumber     uint16
}

type CertificateBaptism struct {
	gorm.Model
	OrgEmisor      uint
	RegID          uint
	Register       RegisterBaptism `gorm:"foreignKey:RegID"`
	UserValidator  uint
	DateEmission   time.Time
	DateExpiration time.Time
	UUID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();<-:create"`
}

type CertificateBaptismPDF struct {
	CertID             string
	CertUUID           string
	CertDateEmission   string
	CertDateExpiration string
	OrgEmiDiocese      string
	OrgEmiType         string
	OrgEmiName         string
	OrgEmiCommune      string
	OrgEmiAddress      string
	RegBookNumber      string
	RegBookPage        string
	RegOrgBaptism      string
	RegDateBaptism     string
	RegUserBaptized    string
	RegUserRUT         string
	RegUserBirthDate   string
	RegUserBirthPlace  string
	RegUserFather      string
	RegUserMother      string
	RegUserGodfather   string
	RegUserGodmother   string
	RegValidator       string
}
