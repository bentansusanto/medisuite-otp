package config

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
)

func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("failed to generate random token")
		return ""
	}
	return hex.EncodeToString(b)
}
