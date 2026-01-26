package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Hash the password with Argon2id and return the hashed password
// salt is generated randomly and appended to the hashed password
// format: "hashedPassword!salt"
func HashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)

	passHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 1, 32)
	passBase64 := base64.StdEncoding.EncodeToString(passHash)

	return fmt.Sprintf("%s!%s", passBase64, saltBase64)
}

// VerifyPassword checks if the provided password matches the stored hashed password
func VerifyPassword(toVerify string, stored string) bool {
	splits := strings.Split(stored, "!")
	salt64 := splits[1]
	storedPass := splits[0]
	saltbytes, err := base64.StdEncoding.DecodeString(salt64)
	if err != nil {
		return false
	}
	toVerifyHash := argon2.IDKey([]byte(toVerify), saltbytes, 1, 64*1024, 1, 32)
	toVerifyHashBase64 := base64.StdEncoding.EncodeToString(toVerifyHash)
	return storedPass == toVerifyHashBase64
}

// random 32 bit hashed string
func GetRandomHash() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	sum := sha256.Sum256(bytes)
	return base64.StdEncoding.EncodeToString(sum[:])
}
