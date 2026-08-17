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

func GetOrganizationFromID(db *gorm.DB, orgID uint) Organization {
	var org Organization
	result := db.First(&org, orgID)
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
func GetMemberFromUUID(db *gorm.DB, memberUUID uuid.UUID) Member {
	var member Member
	result := db.Where("uuid = ?", memberUUID).First(&member)
	if result.Error != nil {
		log.Printf("Error retrieving member by UUID %s: %s", memberUUID, result.Error)
		return member
	}
	log.Printf("N° members in database with UUID %s: %d", memberUUID, result.RowsAffected)
	return member
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

func GetUserFromID(db *gorm.DB, userID uint) User {
	var user User
	result := db.First(&user, userID)
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

func GetMemberFromID(db *gorm.DB, memberID uint) Member {
	var member Member
	result := db.First(&member, memberID)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return member
	}
	log.Printf("N° members in database: %d", result.RowsAffected)
	return member
}
func GetMembers(db *gorm.DB) []Member {
	var members []Member
	result := db.Preload("Organization").Find(&members)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return nil
	}
	log.Printf("N° members in database: %d", result.RowsAffected)
	return members
}
func GetMembersFromOrg(db *gorm.DB, orgID uint) []Member {
	var members []Member
	result := db.Preload("Organization").Where("org_id = ?", orgID).Find(&members)
	if result.Error != nil {
		log.Printf("Error obtaining members: %s", result.Error)
		return nil
	}
	log.Printf("N° members from organization %d in database: %d", orgID, result.RowsAffected)
	return members
}
func GetMembersFromOrgAndRole(db *gorm.DB, orgID uint, role uint8) []Member {
	var members []Member
	result := db.Preload("Organization").Where("org_id = ? AND role = ?", orgID, role).Find(&members)
	if result.Error != nil {
		log.Printf("Error obtaining members: %s", result.Error)
		return nil
	}
	log.Printf("N° members from organization %d with role %d in database: %d", orgID, role, result.RowsAffected)
	return members
}
func createMember(db *gorm.DB, member Member) {
	result := db.Create(&member)
	log.Printf("Member ID: %d", member.ID)
	if result.Error != nil {
		log.Printf("Error inserting member: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
}

func getBookFromID(db *gorm.DB, bookID uint) Book {
	var book Book
	result := db.First(&book, bookID)
	if result.Error != nil {
		log.Printf("Error obtaining book: %s", result.Error)
		return book
	}
	log.Printf("N° books with id: %d in database: %d", bookID, result.RowsAffected)
	return book
}
func getBooksFromOrg(db *gorm.DB, orgID uint) []Book {
	var books []Book
	result := db.Order("book_nr").Where("org_id = ?", orgID).Find(&books)
	if result.Error != nil {
		log.Printf("Error obtaining books: %s", result.Error)
		return nil
	}
	log.Printf("N° books from organization %d in database: %d", orgID, result.RowsAffected)
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

func getIndexBaptismFromID(db *gorm.DB, indexID uint) IndexBaptism {
	var index IndexBaptism
	result := db.First(&index, indexID)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return index
	}
	log.Printf("N° baptism index with id: %d in database: %d", indexID, result.RowsAffected)
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
	regID := regBaptism.ID
	if result.Error != nil {
		log.Printf("Error inserting baptism register: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
	return regID
}
func getRegisterBaptismFromID(db *gorm.DB, regID uint) RegisterBaptism {
	var reg RegisterBaptism
	result := db.First(&reg, regID)
	if result.Error != nil {
		log.Printf("Error retrieving data: %s", result.Error)
		return reg
	}
	log.Printf("N° baptism registers in database: %d", result.RowsAffected)
	return reg
}
func getRegistersBaptismFromBook(db *gorm.DB, bookID uint) []RegisterBaptism {
	var regs []RegisterBaptism
	result := db.Where("book_id = ?", bookID).Find(&regs)
	if result.Error != nil {
		log.Printf("Error obtaining baptism registers: %s", result.Error)
		return nil
	}
	log.Printf("N° pages from book %d in database: %d", bookID, result.RowsAffected)
	return regs
}
func updateOrInsertIndexBaptism(db *gorm.DB, indexID uint, regNew RegisterBaptism) {
	index := getIndexBaptismFromID(db, indexID)
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
	certID := certBaptism.ID
	if result.Error != nil {
		log.Printf("Error inserting baptism certificate: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
	return certID
}
func getCertificatesBaptismFromOrg(db *gorm.DB, orgID uint) []CertificateBaptism {
	var certificates []CertificateBaptism
	result := db.Preload("Register").Where("org_emisor = ?", orgID).Find(&certificates)
	if result.Error != nil {
		log.Printf("Error obtaining baptism certificate: %s", result.Error)
		return nil
	}
	log.Printf("N° certificates from organization %d in database: %d", orgID, result.RowsAffected)
	return certificates
}
func getCertificatesBaptismFromOrgAndReg(db *gorm.DB, orgID uint, regID uint) []CertificateBaptism {
	var certificates []CertificateBaptism
	result := db.Preload("Register").Where("org_emisor = ? AND reg_id = ?", orgID, regID).Find(&certificates)
	if result.Error != nil {
		log.Printf("Error obtaining baptism certificate: %s", result.Error)
		return nil
	}
	log.Printf("N° certificates from register %d emitted by organization %d in database: %d", regID, orgID, result.RowsAffected)
	return certificates
}

func createCertificateRequest(db *gorm.DB, certReq CertificateRequest) uint {
	result := db.Create(&certReq)
	if result.Error != nil {
		log.Printf("Error inserting certificate request: %s", result.Error)
	}
	log.Printf("N° rows inserted: %d", result.RowsAffected)
	return certReq.ID
}
func getCertificateRequestsFromMail(db *gorm.DB, mail string) []CertificateRequest {
	var requests []CertificateRequest
	result := db.Preload("Organization").Where("mail = ?", mail).Find(&requests)
	if result.Error != nil {
		log.Printf("Error obtaining certificate requests from mail: %s", result.Error)
		return nil
	}
	log.Printf("N° certificate requests with mail %s in database: %d", mail, result.RowsAffected)
	return requests
}
func getCertificateRequestsFromOrg(db *gorm.DB, orgID uint) []CertificateRequest {
	var requests []CertificateRequest
	result := db.Preload("Organization").Find(&requests, CertificateRequest{OrgID: orgID})
	if result.Error != nil {
		log.Printf("Error obtaining certificate requests from org: %s", result.Error)
		return nil
	}
	log.Printf("N° certificate requests from organization %d in database: %d", orgID, result.RowsAffected)
	return requests
}

func getBooksIndexBaptismFromOrg(db *gorm.DB, orgID uint) []Book {
	var booksIndex []Book
	result := db.Preload("IndexBaptism", func(tx *gorm.DB) *gorm.DB { return tx.Limit(30) }).Find(&booksIndex, Book{OrgID: orgID})
	if result.Error != nil {
		log.Printf("Error obtaining books with baptism index: %s", result.Error)
		return nil
	}
	log.Printf("N° books with baptism index from org %d in database: %d", orgID, result.RowsAffected)
	return booksIndex
}
