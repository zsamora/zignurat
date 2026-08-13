package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashed)
}

func VerifyPassword(accountPassword, formPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(accountPassword), []byte(formPassword))
	if err != nil {
		log.Printf("Error: %s", err)
		return false
	} else {
		return true
	}
}
