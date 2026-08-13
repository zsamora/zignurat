package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/zsamora/utils"
)

func connectGORM() *gorm.DB {
	host := utils.GetConfig("IDENTIRAT_DB_HOST")
	port := utils.GetConfig("IDENTIRAT_DB_PORT")
	dbname := utils.GetConfig("IDENTIRAT_DB_NAME")
	user := utils.GetConfig("IDENTIRAT_DB_USER")
	password := utils.GetConfig("IDENTIRAT_DB_PW")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed connection to database")
	}
	return db
}

func getAccounts(db *gorm.DB) []Account {
	var accounts []Account
	result := db.Find(&accounts)
	if result.Error != nil {
		log.Printf("# Identirat: Error obtaining accounts: %s", result.Error)
		return nil
	}
	log.Printf("# Identirat: N° accounts in database: %d", result.RowsAffected)
	return accounts
}

func getAccountByUsername(db *gorm.DB, username string) Account {
	var account Account
	result := db.Where("username = ?", username).First(&account)
	if result.Error != nil {
		log.Printf("# Identirat: Error obtaining account with username %s: %s", username, result.Error)
		return account
	}
	log.Printf("# Identirat: N° accounts with username: %s in database: %d", username, result.RowsAffected)
	return account
}

func createAccount(db *gorm.DB, acc Account) {
	result := db.Create(&acc)
	if result.Error != nil {
		log.Printf("# Identirat: Error inserting account: %s", result.Error)
	}
	log.Printf("# Identirat: N° rows inserted: %d", result.RowsAffected)
}

func getAccountById(db *gorm.DB, id int64) Account {
	var account Account
	result := db.First(&account, id)
	if result.Error != nil {
		log.Printf("# Identirat: Error obtaining account with id %d: %s", id, result.Error)
		return account
	}
	log.Printf("# Identirat: N° accounts with id: %d in database: %d", id, result.RowsAffected)
	return account
}

func updateAccount(db *gorm.DB, acc Account) {
	result := db.Save(&acc)
	if result.Error != nil {
		log.Printf("# Identirat: Error updating account: %s", result.Error)
	}
	log.Printf("# Identirat: N° rows updated: %d", result.RowsAffected)
}
