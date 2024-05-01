package utils

import (
	
	"strings"
	"math/rand/v2"
)
func randomString(length int)string{
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i:= 0; i < length; i++{
		sb.WriteByte(characters[rand.IntN(len(characters))])
	}
	return sb.String()
}

func randomNum(min, max int)int {
	return rand.IntN(max-min+1)+min
}

func RandomOwner()string{
	return randomString(6)
}

func RandomBalance() int{
	return randomNum(100, 1000)
}

func RandomCurrency()string{
	currency := []string{"INR","USD", "EUR"}
	return currency[rand.IntN(len(currency))]
}