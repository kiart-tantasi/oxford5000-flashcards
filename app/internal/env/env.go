package env

import "os"

func GetEnv(name, defaultValue string) string {
	val := os.Getenv(name)
	if val != "" {
		return val
	}
	return defaultValue
}
