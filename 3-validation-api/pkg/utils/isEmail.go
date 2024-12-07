package utils

import (
	"strings"
)

func IsEmail(email string, payload string) bool {
	return strings.Contains(email, payload)
}
