package service

import "math/rand"

var quotes = []string{
	"Stay hungry, stay foolish.",
	"Simplicity is the ultimate sophistication.",
	"First solve the problem, then write the code.",
}

func RandomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}