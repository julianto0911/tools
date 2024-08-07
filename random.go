package tools

import (
	"math/rand"

	"github.com/google/uuid"
)

func RandomNumber(min, max int) int {
	if max-min == 0 {
		return 0
	}

	randomNumberPin := rand.Intn(max-min) + min
	return randomNumberPin
}

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyz")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CreateUUID() string {
	uid := uuid.New()
	return uid.String()
}
