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

var sacramentTypes = map[uint8]string{
	0: "Bautismo",
	1: "Confirmación",
	2: "Matrimonio",
}

var certificatePurposes = map[uint8]string{
	0: "Inscripción Colegio",
	1: "Otro",
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

func ParseSacrType(sacrTypeID uint8) string {
	sacrName, found := sacramentTypes[sacrTypeID]
	if found {
		return sacrName
	}
	log.Printf("Key %d not found in sacrament types", sacrTypeID)
	return ""
}

func ParseCertPurpose(certPurposeID uint8) string {
	purposeName, found := certificatePurposes[certPurposeID]
	if found {
		return purposeName
	}
	log.Printf("Key %d not found in certificate purposes", certPurposeID)
	return ""
}
