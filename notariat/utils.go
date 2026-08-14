package main

import (
	"log"
)

var orgTypes = map[uint8]string{
	1: "Parroquia",
	2: "Capilla",
}

var bookTypes = map[uint8]string{
	1: "Índice",
	2: "Bautismo",
}

func ParseOrgType(orgTypeId uint8) string {
	orgName, found := orgTypes[orgTypeId]
	if found {
		return orgName
	}
	log.Printf("Key %d not found in org types", orgTypeId)
	return ""
}

func ParseBookType(bookTypeId uint8) string {
	bookName, found := bookTypes[bookTypeId]
	if found {
		return bookName
	}
	log.Printf("Key %d not found in book types", bookTypeId)
	return ""
}
