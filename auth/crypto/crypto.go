// Package crypto provides support for encrypting passwords and generating
// tokens for authentication support.
package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// SecureEntity interface is required for processing and validating tokens.
type SecureEntity interface {
	Pwd() ([]byte, error)
	Salt() ([]byte, error)
}

// BcryptPassword uses bcrypt to hash password.
func BcryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CompareBcryptHashPassword compares pwd hash to original hash.
func CompareBcryptHashPassword(hash []byte, pwd []byte) error {
	return bcrypt.CompareHashAndPassword(hash, pwd)
}

// SignedHash generates a Signed SHA256 Hash.
func SignedHash(pwd []byte, salt []byte) ([]byte, error) {
	key, err := scrypt.Key([]byte(pwd), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	// Append a salt to the password.
	h := hmac.New(sha256.New, key)
	h.Write(pwd)
	return h.Sum(nil), nil
}

// GenerateToken returns a hash for SecureEntity interface.
func GenerateToken(entity SecureEntity) ([]byte, error) {
	pwd, err := entity.Pwd()
	if err != nil {
		return nil, err
	}

	salt, err := entity.Salt()
	if err != nil {
		return nil, err
	}

	return SignedHash(pwd, salt)
}

// IsTokenValid checks whether a hash is valid.
func IsTokenValid(entity SecureEntity, hash string) error {
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}

	eHash, hErr := GenerateToken(entity)
	if hErr != nil {
		return hErr
	}

	if hmac.Equal(decodedHash, eHash) == false {
		return errors.New("Invalid Token")
	}

	return nil
}
