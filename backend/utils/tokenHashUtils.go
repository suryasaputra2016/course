package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Generate Token takes length of byte as integer and returns token as string
func GenerateToken(lengthByte int) (string, error) {
	tokenByte := make([]byte, lengthByte)
	totalRead, err := rand.Read(tokenByte)
	if err != nil {
		return "", fmt.Errorf("creating random byte: %s", err)
	}
	if totalRead < lengthByte {
		return "", fmt.Errorf("not enough read bytes: %s", err)
	}
	return base64.URLEncoding.EncodeToString(tokenByte), nil
}

// HashToken takes token string and returns token hash string
func HashToken(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
