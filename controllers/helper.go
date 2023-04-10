package controllers

import (
	"log"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// hashPasswd returns the encrpyted password string to store in the database
func hashPasswd(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		log.Println("Error hashing password")
		return "", err
	}

	return string(hash), nil
}

// varifyPasswd verifies the password and returns a boolean value(true if the password is valid).
func verifyPasswd(hashed, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	valid := true
	msg := ""
	if err != nil {
		msg = "Incorrect or invalid password"
		valid = false
		return valid, msg
	}

	return valid, msg
}

// validate simplies initializes an instance from the validator package
func validate() *validator.Validate {
	return validator.New()
}
