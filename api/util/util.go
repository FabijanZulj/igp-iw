package util

import (
	"fmt"
	"math/rand"
	"time"
)

func BuildVerificationUrl(baseUrl, email, verifyCode string) string {
	return fmt.Sprintf("%s/verify/%s/%s", baseUrl, email, verifyCode)
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}
