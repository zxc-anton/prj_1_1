package auth

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func comparePasswords(password, hashPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hashPassword)) != nil
}
