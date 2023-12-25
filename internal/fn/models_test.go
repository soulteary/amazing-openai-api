package fn

import (
	"reflect"
	"testing"

	"github.com/soulteary/amazing-openai-api/internal/define"
)

func TestExtractModelAlias(t *testing.T) {
	var result define.ModelAlias

	tests := []struct {
		name     string
		alias    string
		expected define.ModelAlias
	}{
		{
			name:     "empty string",
			alias:    "",
			expected: result,
		},
		{
			name:     "single valid alias pair",
			alias:    "key1:value1",
			expected: define.ModelAlias{{"key1", "value1"}},
		},
		{
			name:     "multiple valid alias pairs",
			alias:    "key1:value1,key2:value2",
			expected: define.ModelAlias{{"key1", "value1"}, {"key2", "value2"}},
		},
		{
			name:     "invalid alias pair",
			alias:    "singleword",
			expected: result,
		},
		{
			name:     "mixed valid and invalid alias pairs",
			alias:    "key1:value1,singleword,key2:value2",
			expected: define.ModelAlias{{"key1", "value1"}, {"key2", "value2"}},
		},
		{
			name:     "valid alias with extra colon",
			alias:    "key1:value1:extra",
			expected: result,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractModelAlias(tt.alias)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractModelAlias(%q) = %v, want %v", tt.alias, result, tt.expected)
			}
		})
	}
}
