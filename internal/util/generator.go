package util

import (
	"math/rand"
	"time"
)

func GenerateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var charsets = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[rand.Intn(len(charsets))]
	}
	return string(letters)
}

func GenerateRandomNumber(n int) string {
	rand.Seed(time.Now().UnixNano())
	var charsets = []rune("0123456789")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[rand.Intn(len(charsets))]
	}
	return string(letters)
}
