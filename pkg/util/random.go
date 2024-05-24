package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet string = "abcdefghijklmnopqrstuvywxyz"
)

var domainNames = []string{"@gmail.com", "@hotmail.com", "@outlook.com"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomArrayInt(min, max, n int64) []int64 {
	randomNums := make([]int64, 0, n)
	for i := int64(0); i < n; i++ {
		randomNums = append(randomNums, RandomInt(min, max))
	}
	return randomNums
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUsername() string {
	return RandomString(6)
}

func RandomEmail() string {
	return RandomUsername() + RandomDomainName()
}
func RandomDomainName() string {
	n := len(domainNames)
	return domainNames[rand.Intn(n)]
}
