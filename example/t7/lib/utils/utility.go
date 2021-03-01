package utils

import (
	"crypto/sha256"
	"fmt"
	mrand "math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Util ...
type Util struct{}

var (
	letters   = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lcletters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
)

func init() {
	mrand.Seed(time.Now().UnixNano())
}

// Random ...
func Random(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

// RandomLC ...
func RandomLC(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = lcletters[mrand.Intn(len(lcletters))]
	}
	return string(b)
}

// GetHash256 ...
func GetHash256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// PasswordHash ...
func PasswordHash(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

// PasswordVerify ...
func PasswordVerify(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
