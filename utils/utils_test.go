package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	key := "TEST_ENV_1"
	expected := "123456"
	t.Setenv(key, expected)

	value := GetEnvOrDefault(key, "default")
	assert.Equal(t, expected, value)
}

func TestGetEnvOrDefaultForDefault(t *testing.T) {
	expected := "default"

	value := GetEnvOrDefault("TEST_ENV_1", "default")
	assert.Equal(t, expected, value)
}

func TestGetEnvOrFail(t *testing.T) {
	key := "TEST_ENV_1"
	expected := "123456"
	t.Setenv(key, expected)

	value := GetEnvOrFail(key)
	assert.Equal(t, expected, value)
}
