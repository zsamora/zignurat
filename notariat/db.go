package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/zsamora/utils"
)

func connectGORM() *gorm.DB {
	host := utils.GetConfig("NOTARIAT_DB_HOST")
	port := utils.GetConfig("NOTARIAT_DB_PORT")
	dbname := utils.GetConfig("NOTARIAT_DB_NAME")
	user := utils.GetConfig("NOTARIAT_DB_USER")
	password := utils.GetConfig("NOTARIAT_DB_PW")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed connection to database Notariat")
	}
	return db
}

func GetOrganizationFromId(db *gorm.DB, orgId uint) Organization {
	var org Organization
	result := db.First(&org, orgId)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return org
	}
	log.Printf("N° organizations in database: %d", result.RowsAffected)
	return org
}
func GetOrganizationFromUUID(db *gorm.DB, orgUUID uuid.UUID) Organization {
	var org Organization
	result := db.Where("uuid = ?", orgUUID).First(&org)
	if result.Error != nil {
		log.Printf("Error retrieving organization by UUID %s: %s", orgUUID, result.Error)
		return org
	}
	log.Printf("N° organizations in database with UUID %s: %d", orgUUID, result.RowsAffected)
	return org
}
func GetOrganizations(db *gorm.DB) []Organization {
	var orgs []Organization
	result := db.Find(&orgs)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return nil
	}
	log.Printf("N° organizations in database: %d", result.RowsAffected)
	return orgs
}
func createOrganization(db *gorm.DB, org Organization) {
	result := db.Create(&org)
	log.Printf("Organization ID: %d", org.ID)
	if result.Error != nil {
		log.Printf("Error inserting organization: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
}

func GetUserFromId(db *gorm.DB, userId uint) User {
	var user User
	result := db.First(&user, userId)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return user
	}
	log.Printf("N° users in database: %d", result.RowsAffected)
	return user
}
func GetUsers(db *gorm.DB) []User {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return nil
	}
	log.Printf("N° users in database: %d", result.RowsAffected)
	return users
}
func createUser(db *gorm.DB, user User) {
	result := db.Create(&user)
	log.Printf("User ID: %s", user.ID)
	if result.Error != nil {
		log.Printf("Error inserting user: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
}

func getBookFromId(db *gorm.DB, bookId uint) Book {
	var book Book
	result := db.First(&book, bookId)
	if result.Error != nil {
		log.Printf("Error obtaining book: %s", result.Error)
		return book
	}
	log.Printf("N° books with id: %d in database: %d", bookId, result.RowsAffected)
	return book
}
func getBooksFromOrg(db *gorm.DB, orgId uint) []Book {
	var books []Book
	result := db.Order("book_nr").Find(&books, Book{OrgID: orgId})
	if result.Error != nil {
		log.Printf("Error obtaining books: %s", result.Error)
		return nil
	}
	log.Printf("N° books from organization %d in database: %d", orgId, result.RowsAffected)
	return books
}
func createBook(db *gorm.DB, book Book) {
	result := db.Create(&book)
	log.Printf("Book ID: %d", book.ID)
	if result.Error != nil {
		log.Printf("Error inserting book: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
}

func getIndexBaptismFromId(db *gorm.DB, indexId uint) IndexBaptism {
	var index IndexBaptism
	result := db.First(&index, indexId)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return index
	}
	log.Printf("N° baptism index with id: %d in database: %d", indexId, result.RowsAffected)
	return index
}
func getLastIndexBaptismID(db *gorm.DB) uint {
	var index IndexBaptism
	result := db.Last(&index)
	if result.Error != nil {
		log.Printf("Error obtaining last baptism index: %s", result.Error)
	}
	log.Printf("N° index with id: %d in database: %d", index.ID, result.RowsAffected)
	return index.ID
}

func createRegisterBaptism(db *gorm.DB, regBaptism RegisterBaptism) uint {
	result := db.Create(&regBaptism)
	regId := regBaptism.ID
	if result.Error != nil {
		log.Printf("Error inserting baptism register: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
	return regId
}
func getRegisterBaptismFromId(db *gorm.DB, regId uint) RegisterBaptism {
	var reg RegisterBaptism
	result := db.First(&reg, regId)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return reg
	}
	log.Printf("N° baptism registers in database: %d", result.RowsAffected)
	return reg
}
func getRegistersBaptismFromBook(db *gorm.DB, bookId uint) []RegisterBaptism {
	var regs []RegisterBaptism
	result := db.Find(&regs, RegisterBaptism{BookID: bookId})
	if result.Error != nil {
		log.Printf("Error obtaining baptism registers: %s", result.Error)
		return nil
	}
	log.Printf("N° pages from book %d in database: %d", bookId, result.RowsAffected)
	return regs
}
func updateOrInsertIndexBaptism(db *gorm.DB, indexId uint, regNew RegisterBaptism) {
	index := getIndexBaptismFromId(db, indexId)
	index.OrgID = regNew.OrgBaptism
	index.BookID = regNew.BookID
	index.RegID = regNew.ID
	index.PageNumber = regNew.PageNumber
	index.UserSurnameF = regNew.FatherSurname
	index.UserSurnameM = regNew.MotherSurname
	index.UserNameFirst = regNew.BaptizedNameF
	index.UserNameSecond = regNew.BaptizedNameS
	result := db.Save(&index)
	if result.Error != nil {
		log.Printf("Error obtaining index: %s", result.Error)
	}
	log.Printf("N° rows inserted or updated: %d", result.RowsAffected)
}

func createCertificateBaptism(db *gorm.DB, certBaptism CertificateBaptism) uint {
	result := db.Create(&certBaptism)
	certId := certBaptism.ID
	if result.Error != nil {
		log.Printf("Error inserting baptism certificate: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
	return certId
}
func getCertificatesBaptismFromOrg(db *gorm.DB, orgId uint) []CertificateBaptism {
	var certificates []CertificateBaptism
	result := db.Preload("Register").Find(&certificates, CertificateBaptism{OrgEmisor: orgId})
	if result.Error != nil {
		log.Printf("Error obtaining baptism certificate: %s", result.Error)
		return nil
	}
	log.Printf("N° certificates from organization %d in database: %d", orgId, result.RowsAffected)
	return certificates
}
func getCertificatesBaptismFromOrgAndReg(db *gorm.DB, orgId uint, regId uint) []CertificateBaptism {
	var certificates []CertificateBaptism
	result := db.Preload("Register").Find(&certificates, CertificateBaptism{OrgEmisor: orgId, RegID: regId})
	if result.Error != nil {
		log.Printf("Error obtaining baptism certificate: %s", result.Error)
		return nil
	}
	log.Printf("N° certificates from register %d emitted by organization %d in database: %d", regId, orgId, result.RowsAffected)
	return certificates
}
