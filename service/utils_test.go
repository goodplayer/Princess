package service

import (
	"bytes"
	"testing"
)

func TestUtf8DataLength(t *testing.T) {
	raw := "哈哈哈😄"
	t.Log(len([]byte(raw)) == 13, len(bytes.Runes([]byte(raw))) == 4)
}

func TestPasswordEncrypt(t *testing.T) {
	password := "aaaa"
	output := PasswordEncrypt(password)
	t.Log(output, len(output) == 64)
}
