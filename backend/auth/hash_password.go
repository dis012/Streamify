package auth

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func CheckIfPasswordValid(password string) bool {
	/*
		Check if the password is correct length and that it contains at least one uppercase letter and
		at least one symbol (special character)
	*/
	hasUpper := false
	hasSymbol := false
	correctLen := false

	if len(password) >= 8 {
		correctLen = true
	}

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}

		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			hasSymbol = true
		}

		if hasUpper && hasSymbol && correctLen {
			return true
		}
	}

	return false
}

func HashPassword(password string) (string, error) {
	/*
		Function that hashes the given password using bcrypt.
		It accepts the password as a string and returns the hashed password as a string.
	*/

	checkPassword := CheckIfPasswordValid(password)
	if !checkPassword {
		return "", errors.New("password must contain at least one uppercase letter, at least one special character and be at least 8 char long")
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashed_password), err
}

func CompareSavedAndInputPassword(password, hashed_password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
	return err == nil
}
