package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

func MimcHash(data []byte) string {
	f := bn254.NewMiMC()
	f.Write(data)
	hash := f.Sum(nil)
	if len(hash) < 32 {
		padding := make([]byte, 32-len(hash))
		hash = append(hash, padding...)
	}
	hashInt := big.NewInt(0).SetBytes(hash)
	return hashInt.String()
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func SignHMAC256(msg, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return hex.EncodeToString(mac.Sum(nil))
}
