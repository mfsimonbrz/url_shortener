package utils

import (
	"math/rand/v2"
	"strings"
)

const (
	SHORT_URL_SIZE      = 11
	SHORT_URL_GENERATOR = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateRandomString() string {
	var sb strings.Builder

	for i := 0; i < SHORT_URL_SIZE; i++ {
		pos := rand.IntN(len(SHORT_URL_GENERATOR) - 1)
		sb.WriteString(string(SHORT_URL_GENERATOR[pos]))
	}

	return sb.String()
}
