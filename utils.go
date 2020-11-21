package main

import (
	"fmt"
	"hash/fnv"
	"os"
)

// Get env
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Generate code from URL
func genCode(url string) string {
	h := fnv.New32a()
	h.Write([]byte(url))
	return fmt.Sprint(h.Sum32())
}
