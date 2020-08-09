package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func EncryptScrypt(passwordPlainText string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return encryptInternal([]byte(passwordPlainText), salt)
}

func encryptInternal(passwordPlainText, salt []byte) (string, error) {
	pw, err := scrypt.Key(passwordPlainText, salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(base64.StdEncoding.EncodeToString(pw)) + "." + fmt.Sprint(base64.StdEncoding.EncodeToString(salt)), nil
}

func VerifyScrypt(password string, passwordPlainText string) (bool, error) {
	idx := strings.IndexAny(password, ".")
	if idx == -1 {
		return false, errors.New("cannot find dot char when decrypting scrypt password")
	}
	if idx+1 == len(password) {
		return false, errors.New("no salt found when decrypting scrypt password")
	}
	salt, err := base64.StdEncoding.DecodeString(password[idx+1:])
	if err != nil {
		return false, err
	}
	pw, err := encryptInternal([]byte(passwordPlainText), salt)
	if err != nil {
		return false, err
	}
	if subtle.ConstantTimeCompare([]byte(pw), []byte(password)) == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
