package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

func PasswordEncrypt(rawPassword string) string {
	data := sha256.Sum256([]byte(rawPassword))
	return hex.EncodeToString(data[:])
}

func CheckUsernameLength(username string) bool {
	l := len(bytes.Runes([]byte(username)))
	return l > 0 && l <= 20
}

func CheckPasswordLength(password string) bool {
	l := len([]byte(password))
	return l > 0 && l <= 64
}

func CheckNicknameLength(nickname string) bool {
	l := len(bytes.Runes([]byte(nickname)))
	return l > 0 && l <= 20
}

func CheckEmailLength(email string) bool {
	l := len([]byte(email))
	return l > 0 && l <= 50
}
