package pelican

import (
	"math/rand"
	"time"
)

// letterBytes are the letters used to generate a random string.
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// randomString generates a string consisting of characters from letterBytes that is n characters long.
// Courtesy: https://stackoverflow.com/a/31832326
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func init() {
	// Seed the random number generator.
	rand.Seed(time.Now().UnixNano())
}
