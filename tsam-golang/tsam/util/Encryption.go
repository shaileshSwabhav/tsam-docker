package util

import (
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword will encrypt the password and return encrypted string
func EncryptPassword(password string) (string, error) {

	cost := 8
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return password, err
	}

	// login.Password = string(hashPassword)
	return string(hashPassword), nil
}

// DoPasswordsMatch will compare the given loginPassword with the encrypted password in the db
func DoPasswordsMatch(hashedPassword, password string) bool {

	// password, err := bcrypt.GenerateFromPassword([]byte(loginPassword), 8)
	// if err != nil {
	// 	return errors.NewValidationError("GenerateFromPassword Error ")
	// }
	// fmt.Println("===================================================")
	// // 9@@VK4pQ
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}
