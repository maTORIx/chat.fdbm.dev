package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(s, salt string) string {
	r := sha256.Sum256([]byte(s + salt))
	return hex.EncodeToString(r[:])
}

func GenerateUserID(ipAddress, salt string) string {
	return GenerateHash(ipAddress, salt)
}
