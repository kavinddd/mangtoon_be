package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomInt(t *testing.T) {

	require.Equal(t, int64(0), RandomInt(0, 0))
	require.NotZero(t, RandomInt(1, 999))

}

func TestRandomArrayInt(t *testing.T) {
	tenRandomNums := RandomArrayInt(1, 5, 10)
	require.Len(t, tenRandomNums, 10)

	isValid := true
	randomNums := RandomArrayInt(50, 51, 100)
	for _, num := range randomNums {
		if num != 50 && num != 51 {
			isValid = false
			break
		}
	}
	require.True(t, isValid)
}

func TestRandomDomainName(t *testing.T) {

}

func TestRandomString(t *testing.T) {

}

func TestRandomUsername(t *testing.T) {

}

func TestRandomEmail(t *testing.T) {

}
