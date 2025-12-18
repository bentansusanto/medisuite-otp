package config

import (
	"log/slog"

	errWrap "medisuite-api/common/errors"
	errConstant "medisuite-api/constants/errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// generate password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// check if error hashing password
	if err != nil {
		slog.Error("Error hashing password")
		return "", errWrap.WrapError(errConstant.ErrSQLError)
	}
	// convert to string and return
	return string(bytes), nil
}

// verify password
func VerifyPassword(password string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Error("Error verifying password")
		return false, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return true, nil
}
