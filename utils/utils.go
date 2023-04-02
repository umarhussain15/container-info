package utils

import (
	"log"
	"os"
)

// GetEnvOrDefault try to get the env variable value if present, otherwise return default value.
func GetEnvOrDefault(key, defaultValue string) string {
	env, present := os.LookupEnv(key)
	if !present {
		return defaultValue
	}
	return env
}

// GetEnvOrFail exit the program with error if given env variable is not present.
func GetEnvOrFail(key string) string {
	env, present := os.LookupEnv(key)
	if !present {
		log.Fatalln("cannot find the env variable:", key)
	}
	return env
}
