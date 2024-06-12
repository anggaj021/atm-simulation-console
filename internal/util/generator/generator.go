package generator

import (
	"math/rand"
)

func GenerateRandomNDigitNumber(n int) int {
	if n <= 0 {
		return 0 // Invalid input
	}

	// Calculate the range for the n-digit number
	min := 1
	for i := 1; i < n; i++ {
		min *= 10
	}
	max := min*10 - 1

	// Generate a random number within the specified range
	return rand.Intn(max-min+1) + min
}
