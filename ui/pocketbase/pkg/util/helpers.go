package util

import (
	"math/rand"
	"strings"
)

var (
	SpecialCharacters string = "!@#$%^&*()_+-=[]{}|;':,./<>?`~"
	Characters               = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
)

func StringParser(s string) string {
	returnString := strings.ToLower(s)
	returnString = strings.Replace(returnString, " ", "-", -1)
	// remove special characters from returnString
	for _, c := range SpecialCharacters {
		returnString = strings.Replace(returnString, string(c), "", -1)
	}
	return returnString
}

func GenerateId(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = Characters[rand.Intn(len(Characters))]
	}
	return string(b)
}

func StringContains(s string, substr string) bool {
	return strings.Contains(s, substr)
}
