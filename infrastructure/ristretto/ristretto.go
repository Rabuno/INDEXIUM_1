package ristretto

import (
	"Test2/internal/domain"
	"log"

	ristretto "github.com/dgraph-io/ristretto/v2"
)

var Cache *ristretto.Cache[string, []domain.Post]
var err error

func InitRistretto() {
	Cache, err = ristretto.NewCache(&ristretto.Config[string, []domain.Post]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		log.Fatalf("Failed to initialize Ristretto cache: %v", err)
	}

	log.Printf("Ristretto cache initialized successfully")
}
