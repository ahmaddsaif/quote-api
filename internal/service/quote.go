package service

import (
	"context"
	"math/rand"
)

type QuoteService struct{}

func NewQuoteService() *QuoteService {
	return &QuoteService{}
}

var quotes = []string{
	"Stay hungry, stay foolish.",
	"Simplicity is the ultimate sophistication.",
	"First solve the problem, then write the code.",
}

func (s *QuoteService) RandomQuote(ctx context.Context) (string, error) {
	return quotes[rand.Intn(len(quotes))], nil
}
