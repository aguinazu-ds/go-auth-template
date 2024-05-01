package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
)

func ContainsUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func ContainsLowercase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

func ContainsNumber(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func ContainsSymbol(s string) bool {
	for _, r := range s {
		if (r >= '!' && r <= '/') || (r >= ':' && r <= '@') || (r >= '[' && r <= '`') || (r >= '{' && r <= '~') {
			return true
		}
	}
	return false
}

func ContainsNoWhitespace(s string) bool {
	for _, r := range s {
		if r == ' ' {
			return false
		}
	}
	return true
}

func GenerateRandomTokenHash() ([]byte, string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, "", err
	}
	plainText := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plainText))
	return hash[:], plainText, nil
}
