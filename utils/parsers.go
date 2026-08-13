package utils

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func FormatDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%02d/%d", day, month, year)
}
func FormatDatetime(t time.Time) string {
	return fmt.Sprintf(t.Format("02/01/2006 15:04:05"))
}

func FormatYear(t time.Time) string {
	year := t.Year()
	return fmt.Sprintf("%d", year)
}
func FormatUint(val uint64) string {
	return strconv.FormatUint(val, 10)
}

func ParseInt(intString string) (int, error) {
	return strconv.Atoi(intString)
}
func ParseUUID(idString string) uuid.UUID {
	idUUID, err := uuid.Parse(idString)
	if err != nil {
		log.Printf("Error parsing UUID: %s", err)
	}
	return idUUID
}
func ParseYearForm(dateString string) time.Time {
	const layout = "2006" // Standard layout date (https://go.dev/src/time/format.go)
	dateResult, err := time.Parse(layout, dateString)
	if err != nil {
		log.Printf("Error parsing date: %s, using 0 as year", err)
		dateResult, err = time.Parse(layout, "0")
	}
	return dateResult
}
func ParseDateForm(dateString string) time.Time {
	const layout = "2006-01-02" // Standard layout date (https://go.dev/src/time/format.go)
	dateResult, err := time.Parse(layout, dateString)
	if err != nil {
		log.Printf("Error parsing date: %s", err)
	}
	return dateResult
}
