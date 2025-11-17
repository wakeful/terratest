package formatting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBackendConfigAsArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  map[string]interface{}
		expect []string
	}{
		{
			name:   "empty config",
			input:  map[string]interface{}{},
			expect: []string{},
		},
		{
			name:   "string value",
			input:  map[string]interface{}{"bucket": "my-bucket"},
			expect: []string{"-backend-config=bucket=my-bucket"},
		},
		{
			name:   "nil value omitted",
			input:  map[string]interface{}{"key": nil},
			expect: []string{"-backend-config=key"},
		},
		{
			name:   "multiple values",
			input:  map[string]interface{}{"region": "us-east-1", "bucket": "state"},
			expect: []string{"-backend-config=bucket=state", "-backend-config=region=us-east-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBackendConfigAsArgs(tt.input)
			assert.ElementsMatch(t, tt.expect, result)
		})
	}
}

func TestFormatPluginDirAsArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name:   "empty path",
			input:  "",
			expect: nil,
		},
		{
			name:   "valid path",
			input:  "/path/to/plugins",
			expect: []string{"-plugin-dir=/path/to/plugins"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPluginDirAsArgs(tt.input)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestToHclString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"nil", nil, "null"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"list of strings", []string{"a", "b"}, `["a", "b"]`},
		{"list of ints", []int{1, 2, 3}, "[1, 2, 3]"},
		{"map", map[string]string{"key": "value"}, `{"key" = "value"}`},
		{"nested list", []interface{}{[]int{1, 2}}, "[[1, 2]]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toHclString(tt.input, false)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestToHclStringNested(t *testing.T) {
	t.Parallel()

	// Nested strings should be quoted
	result := toHclString("nested", true)
	assert.Equal(t, `"nested"`, result)

	// Non-nested strings should not be quoted
	result = toHclString("not-nested", false)
	assert.Equal(t, "not-nested", result)
}
