package tools

import "golang.org/x/crypto/bcrypt"

func HashBcrypt(password string, strength int) (string, error) {
	if strength == 0 {
		return password, nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), strength)
	return string(bytes), err
}
