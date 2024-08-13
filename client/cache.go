package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func cacheToken(token string) error {
	cache, err := getCachePath()
	if err != nil {
		return err
	}

	log.Println("caching token")

	return os.WriteFile(cache, []byte(token), 0644)
}

func getCachePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, "Documents", "token.txt"), nil
}

func getCachedToken() (string, error) {
	cache, err := getCachePath()
	data, err := os.ReadFile(cache)

	if err != nil {
		return "", err
	}

	token := string(data)
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	log.Println("using cached token")

	return token, nil
}
