package util

import (
	"crypto/rand"
	"math"
	"math/big"
)

func GenerateRandomString(n int) string {
	var charsets = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return generateFromCharset(n, charsets)
}

func GenerateRandomNumber(n int) string {
	var charsets = []rune("0123456789")
	return generateFromCharset(n, charsets)
}

func generateFromCharset(n int, charsets []rune) string {
	if n <= 0 || len(charsets) == 0 {
		return ""
	}
	letters := make([]rune, n)
	max := big.NewInt(int64(len(charsets)))
	for i := range letters {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			letters[i] = charsets[0]
			continue
		}
		letters[i] = charsets[num.Int64()]
	}
	return string(letters)
}

func RoundCurrency(val float64) float64 {
	return math.Round(val*100) / 100
}
