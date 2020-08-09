package password

import (
	"errors"
	"testing"
)

func TestEncryptScrypt(t *testing.T) {
	v1, _ := EncryptScrypt("1")
	v2, _ := EncryptScrypt("kjsdhfgvijwvh82vhneakgvljawhgv")
	t.Log(v1, len(v1))
	t.Log(v2, len(v2))
}

func TestVerifyScrypt(t *testing.T) {
	pw1, _ := EncryptScrypt("1")
	pw2, _ := EncryptScrypt("2")
	pw3, _ := EncryptScrypt("3")

	if v, _ := VerifyScrypt(pw1, "1"); !v {
		t.Fatal(errors.New("pw1 not match"))
	}
	if v, _ := VerifyScrypt(pw2, "1"); v {
		t.Fatal(errors.New("pw2 should not match"))
	}
	if v, _ := VerifyScrypt(pw3, "1"); v {
		t.Fatal(errors.New("pw3 should not match"))
	}
	if v, _ := VerifyScrypt(pw2, "2"); !v {
		t.Fatal(errors.New("pw2 not match"))
	}
	if v, _ := VerifyScrypt(pw3, "3"); !v {
		t.Fatal(errors.New("pw3 not match"))
	}
}

func BenchmarkVerifyScrypt(b *testing.B) {
	pw1, _ := EncryptScrypt("1")

	b.ResetTimer()
	if v, _ := VerifyScrypt(pw1, "1"); !v {
		b.Fatal(errors.New("pw1 not match"))
	}
}
