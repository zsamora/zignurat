package utils

import (
	"os"
)

func GetConfig(envKey string) string {
	envValue := os.Getenv(envKey)
	return envValue
}
