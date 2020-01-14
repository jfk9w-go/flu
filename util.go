package flu

import (
	"math/rand"
)

var symbols = []rune("" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"1234567890!@#$%^&*-_+=?~:;")

func GenerateID(length int) string {
	id := make([]rune, length)
	for i := 0; i < length; i++ {
		id[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(id)
}
