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

func ParseOrgType(orgTypeID uint8) string {
	orgName, found := orgTypes[orgTypeID]
	if found {
		return orgName
	}
	log.Printf("Key %d not found in org types", orgTypeID)
	return ""
}

func ParseBookType(bookTypeID uint8) string {
	bookName, found := bookTypes[bookTypeID]
	if found {
		return bookName
	}
	log.Printf("Key %d not found in book types", bookTypeID)
	return ""
}
