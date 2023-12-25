package fn_test

import (
	"os"
	"testing"

	"github.com/soulteary/amazing-openai-api/internal/fn"
)

func TestGetIntOrDefaultFromEnv(t *testing.T) {
	const defaultVal = 10
	const envKey = "TEST_INT_ENV_VAR"

	t.Run("ReturnsDefaultValueForUnset", func(t *testing.T) {
		os.Unsetenv(envKey)
		if got := fn.GetIntOrDefaultFromEnv(envKey, defaultVal); got != defaultVal {
			t.Errorf("Expected default value %d, got %d", defaultVal, got)
		}
	})

	t.Run("ReturnsParsedValue", func(t *testing.T) {
		expected := 42
		os.Setenv(envKey, "42")
		defer os.Unsetenv(envKey)
		if got := fn.GetIntOrDefaultFromEnv(envKey, defaultVal); got != expected {
			t.Errorf("Expected parsed value %d, got %d", expected, got)
		}
	})

	t.Run("IgnoresInvalidValue", func(t *testing.T) {
		os.Setenv(envKey, "invalid")
		defer os.Unsetenv(envKey)
		if got := fn.GetIntOrDefaultFromEnv(envKey, defaultVal); got != defaultVal {
			t.Errorf("Expected default value %d when variable is invalid, got %d", defaultVal, got)
		}
	})
}

func TestGetStringOrDefaultFromEnv(t *testing.T) {
	const defaultVal = "default"
	const envKey = "TEST_STRING_ENV_VAR"

	t.Run("ReturnsDefaultValueForUnset", func(t *testing.T) {
		os.Unsetenv(envKey)
		if got := fn.GetStringOrDefaultFromEnv(envKey, defaultVal); got != defaultVal {
			t.Errorf("Expected default value %s, got %s", defaultVal, got)
		}
	})

	t.Run("ReturnsNonEmptyValue", func(t *testing.T) {
		expected := "test value"
		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetStringOrDefaultFromEnv(envKey, defaultVal); got != expected {
			t.Errorf("Expected non-empty value %s, got %s", expected, got)
		}
	})

	t.Run("TrimsWhitespace", func(t *testing.T) {
		expected := "test value"
		os.Setenv(envKey, "  "+expected+"  ")
		defer os.Unsetenv(envKey)
		if got := fn.GetStringOrDefaultFromEnv(envKey, defaultVal); got != expected {
			t.Errorf("Expected trimmed value %s, got %s", expected, got)
		}
	})
}

func TestIsValidIPAddress(t *testing.T) {
	testCases := []struct {
		ip    string
		valid bool
	}{
		{"192.168.1.1", true},
		{"255.255.255.255", true},
		{"0.0.0.0", true},
		{"256.1.1.1", false},
		{"192.168.1", false},
		{"not an ip", false},
		{"::1", true}, // IPv6
	}

	for _, tc := range testCases {
		t.Run(tc.ip, func(t *testing.T) {
			if got := fn.IsValidIPAddress(tc.ip); got != tc.valid {
				t.Errorf("IsValidIPAddress(%q) = %v; want %v", tc.ip, got, tc.valid)
			}
		})
	}
}
