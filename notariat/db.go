package main

import (
	"fmt"
	"log"

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
