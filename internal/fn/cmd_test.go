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

func TestGetBoolOrDefaultFromEnv(t *testing.T) {
	const envKey = "TEST_BOOL_ENV_VAR"

	t.Run("ReturnsDefaultValueForUnset", func(t *testing.T) {
		os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, false); got != false {
			t.Errorf("Expected default value %v, got %v", false, got)
		}

		os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, true); got != true {
			t.Errorf("Expected default value %v, got %v", true, got)
		}
	})

	t.Run("test on", func(t *testing.T) {
		expected := "on"
		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, false); got != true {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}

		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, true); got != true {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}
	})

	t.Run("test true", func(t *testing.T) {
		expected := "true"
		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, false); got != true {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}

		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, true); got != true {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}
	})

	t.Run("test 1", func(t *testing.T) {
		expected := "1"
		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, false); got != true {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}

		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, true); got != true {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}
	})

	t.Run("test yes", func(t *testing.T) {
		expected := "yes"
		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, false); got != true {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}

		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, true); got != true {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}
	})

	t.Run("test 0", func(t *testing.T) {
		expected := "0"
		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, false); got != false {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}

		os.Setenv(envKey, expected)
		defer os.Unsetenv(envKey)
		if got := fn.GetBoolOrDefaultFromEnv(envKey, true); got != false {
			t.Errorf("Expected non-empty value %v, got %v", expected, got)
		}
	})

}
