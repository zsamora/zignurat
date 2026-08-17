package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
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

type ownerOption struct {
	UUID        uuid.UUID `json:"uuid"`
	DisplayName string    `json:"display_name"`
}

func fetchNotariatOwners(path string) []ownerOption {
	notariatUrl := utils.GetConfig("NOTARIAT_URL")
	notariatPort := utils.GetConfig("NOTARIAT_PORT")
	requestURL := fmt.Sprintf("http://%s:%s/internal/%s", notariatUrl, notariatPort, path)
	res, err := http.Get(requestURL)
	if err != nil {
		log.Printf("# Identirat: Failed fetching Notariat %s: %s", path, err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Printf("# Identirat: Error status code fetching Notariat %s: %d", path, res.StatusCode)
		return nil
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("# Identirat: Failed reading Notariat %s response: %s", path, err)
		return nil
	}
	var owners []ownerOption
	if err := json.Unmarshal(body, &owners); err != nil {
		log.Printf("# Identirat: Failed parsing Notariat %s response: %s", path, err)
		return nil
	}
	return owners
}

func ownersForAccount(module uint8, accRole uint8) []ownerOption {
	switch module {
	case 2:
		return ownersForNotariatRole(accRole)
	default:
		return nil
	}
}

type accountRow struct {
	Account
	Owners       []ownerOption
	OwnerUUIDStr string
	ModuleName   string
	RoleName     string
}

func roleNameFor(module uint8, accRole uint8) string {
	switch module {
	case 1:
		return identiratRoles[accRole]
	case 2:
		return notariatRoles[accRole]
	default:
		return ""
	}
}

func buildAccountRows(accounts []Account) []accountRow {
	type roleKey struct {
		module  uint8
		accRole uint8
	}
	ownersByRole := make(map[roleKey][]ownerOption)
	rows := make([]accountRow, 0, len(accounts))
	for _, account := range accounts {
		key := roleKey{account.Module, account.AccRole}
		owners, cached := ownersByRole[key]
		if !cached {
			owners = ownersForAccount(account.Module, account.AccRole)
			ownersByRole[key] = owners
		}
		ownerUUIDStr := ""
		if account.OwnerUUID != nil {
			ownerUUIDStr = account.OwnerUUID.String()
		}
		rows = append(rows, accountRow{
			Account:      account,
			Owners:       owners,
			ModuleName:   modules[account.Module],
			RoleName:     roleNameFor(account.Module, account.AccRole),
			OwnerUUIDStr: ownerUUIDStr,
		})
	}
	return rows
}

func ownersForNotariatRole(accRole uint8) []ownerOption {
	switch accRole {
	case 2:
		return fetchNotariatOwners("organizations")
	case 3:
		return fetchNotariatOwners("members")
	case 4:
		return fetchNotariatOwners("users")
	default:
		return nil
	}
}
