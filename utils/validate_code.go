package utils

import (
	"fmt"
	"math/rand"
)

func GenerateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
