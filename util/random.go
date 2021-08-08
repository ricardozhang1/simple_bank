package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().Unix())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generate a random string and length equal n
func RandomString(n int) string {
	 var sb strings.Builder
	 k := len(alphabet)

	 for i:=0; i<n; i++ {
	 	c := alphabet[rand.Intn(k)]
	 	sb.WriteByte(c)
	 }
	 return sb.String()
}

// RandomOwnerName generate a random owner name
func RandomOwnerName() string {
	return RandomString(6)
}

// RandomMoney generate a random money number
func RandomMoney() int64 {
	return RandomInt(100, 800)
}

// RandomCurrency generate a random currency code
func RandomCurrency() string {
	currencies := []string{
		"USD", "RMB", "EUR", "CAD",
	}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomUserName() string {
	return RandomOwnerName()
}

func RandomHashedPassword() string {
	return RandomString(16)
}

func RandomFullName() string {
	return RandomString(8)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(8))
}


