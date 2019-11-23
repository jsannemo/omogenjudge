// Package strings contains string-related utilities.
package strings

import (
	"math/rand"
	str "strings"
	"time"
)

// RandStr returns a random base64 string with the given length.
func RandStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789-_")
	var b str.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
