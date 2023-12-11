package hsh

import (
	"crypto/sha256"
	"fmt"
)

func Hash(data []byte) ([]byte, error) {
	h := sha256.New()

	// Write to hash. Never returns an error
	_, err := h.Write(data)
	if err != nil {
		return nil, fmt.Errorf("probably earth is broken: %w", err)
	}

	return h.Sum(nil), nil
}
