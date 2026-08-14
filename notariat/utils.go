package main

import (
	"log"
)

var orgTypes = map[uint8]string{
	1: "Parroquia",
	2: "Capilla",
}

func ParseOrgType(orgTypeId uint8) string {
	orgName, found := orgTypes[orgTypeId]
	if found {
		return orgName
	}
	log.Printf("Key %d not found in org types", orgTypeId)
	return ""
}
