package util

import (
	"math/rand"
	"time"
)

func GenerateOtpValue() string {
	rand.Seed(time.Now().UnixNano())
	const digits = "0123456789"
	otpValue := make([]byte, 6)
	for i := range otpValue {
		otpValue[i] = digits[rand.Intn(len(digits))]
	}
	return string(otpValue)
}
