package controllers

import (
	"log"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func hashPasswd(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		log.Println("Error hashing password")
		return "", err
	}

	return string(hash), nil
}

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

func validate() *validator.Validate {
	return validator.New()
}
