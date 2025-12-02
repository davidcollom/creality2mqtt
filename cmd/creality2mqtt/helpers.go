package main

import (
	"os"
	"strconv"
	"time"
)

func getEnvOrDefault(envKey, defaultVal string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return defaultVal
}
func getEnvOrDefaultDuration(envKey string, defaultVal time.Duration) time.Duration {
	{
		if v := os.Getenv(envKey); v != "" {
			// best effort parse
			if n, err := strconv.Atoi(v); err == nil {
				return time.Duration(n) * time.Second
			}
		}
		return defaultVal
	}
}
