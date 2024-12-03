package generator

import (
	"math/rand"
)

func GenerateRandomNDigitNumber(length int) int {
	if length <= 0 {
		return 0 // Invalid input
	}

	// Calculate the range for the n-digit number
	min := 1
	for i := 1; i < length; i++ {
		min *= 10
	}
	max := min*10 - 1

	// Generate a random number within the specified range
	return rand.Intn(max-min+1) + min
}
