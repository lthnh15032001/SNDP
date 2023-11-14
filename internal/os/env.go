package os

import (
	"fmt"
	"os"
)

// GetEnv retrieves the value of an environment variable having the specified
// key. If the value is empty string, a specified default is returned instead.
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// MustGetEnv retrieves the value of an environment variable having the specified key, panic if key not set
func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%v not set", key))
	}
	return value
}
