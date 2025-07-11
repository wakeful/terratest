package terragrunt

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test the basic outputArgs function behavior
func TestTgOutputArgs(t *testing.T) {
	t.Parallel()

	options := &Options{
		TerragruntDir:    "test/fixtures/terragrunt/terragrunt-output",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Test with a specific key
	args := outputArgs(options, "test_key")
	expected := []string{"-no-color", "test_key"}
	assert.Equal(t, expected, args)

	// Test with empty key (get all outputs)
	args = outputArgs(options, "")
	expected = []string{"-no-color"}
	assert.Equal(t, expected, args)
}

// Test outputArgs with various ExtraArgs combinations
func TestTgOutputArgsWithExtraArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		key       string
		extraArgs []string
		expected  []string
	}{
		{
			name:      "Basic output with key",
			key:       "vpc_id",
			extraArgs: []string{},
			expected:  []string{"-no-color", "vpc_id"},
		},
		{
			name:      "JSON output with key",
			key:       "vpc_id",
			extraArgs: []string{"-json"},
			expected:  []string{"-no-color", "-json", "vpc_id"},
		},
		{
			name:      "All outputs as JSON",
			key:       "",
			extraArgs: []string{"-json"},
			expected:  []string{"-no-color", "-json"},
		},
		{
			name:      "Output with state file",
			key:       "vpc_id",
			extraArgs: []string{"-state=terraform.tfstate"},
			expected:  []string{"-no-color", "-state=terraform.tfstate", "vpc_id"},
		},
		{
			name:      "JSON output with multiple flags",
			key:       "vpc_id",
			extraArgs: []string{"-json", "-lock=false", "-lock-timeout=10s"},
			expected:  []string{"-no-color", "-json", "-lock=false", "-lock-timeout=10s", "vpc_id"},
		},
		{
			name:      "Multiple flags without key",
			key:       "",
			extraArgs: []string{"-json", "-lock=false"},
			expected:  []string{"-no-color", "-json", "-lock=false"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := &Options{
				TerragruntDir:    "test/fixtures/terragrunt/terragrunt-output",
				TerragruntBinary: "terragrunt",
				Logger:           logger.Discard,
				ExtraArgs:        tc.extraArgs,
			}

			args := outputArgs(options, tc.key)
			assert.Equal(t, tc.expected, args)
		})
	}
}

// Integration test using actual terragrunt stack fixture
func TestTgOutputIntegration(t *testing.T) {
	t.Parallel()

	// Create a temporary copy of the stack fixture
	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Initialize and apply terragrunt using stack commands
	_, err = TgStackInitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		ExtraArgs:        []string{"apply", "-auto-approve"},
	}
	_, err = TgStackRunE(t, applyOptions)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			ExtraArgs:        []string{"destroy", "-auto-approve"},
		}
		_, _ = TgStackRunE(t, destroyOptions)
	}()

	// Test string stack output - get output from mother unit
	strOutput := TgOutput(t, options, "mother.output")
	assert.Contains(t, strOutput, "mother/test.txt")

	// Test getting stack output as JSON - note that our cleaning function will still extract just the value
	jsonOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		ExtraArgs:        []string{"-json"},
	}

	strOutputJson := TgOutput(t, jsonOptions, "mother.output")
	// The JSON output for a single value should still be cleaned to just show the value
	assert.Contains(t, strOutputJson, "mother/test.txt")

	// Test getting all stack outputs as JSON
	allOutputsJson := TgOutput(t, jsonOptions, "")
	require.NotEmpty(t, allOutputsJson)

	// For JSON output of all outputs, we should get valid JSON
	// But our function cleans it, so let's test it as-is
	// The JSON structure should be valid and contain our expected data
	if strings.Contains(allOutputsJson, "{") {
		// Parse and validate the JSON structure
		var allOutputs map[string]interface{}
		err = json.Unmarshal([]byte(allOutputsJson), &allOutputs)
		require.NoError(t, err)

		// Verify all expected stack outputs are present
		require.Contains(t, allOutputs, "mother.output")
		require.Contains(t, allOutputs, "father.output")
		require.Contains(t, allOutputs, "chick_1.output")
		require.Contains(t, allOutputs, "chick_2.output")

		// Verify the structure of outputs (should have "value" field)
		motherOutputMap := allOutputs["mother.output"].(map[string]interface{})
		assert.Contains(t, motherOutputMap["value"], "mother/test.txt")
	} else {
		// If not JSON format, at least verify it contains our expected values
		assert.Contains(t, allOutputsJson, "mother.output")
		assert.Contains(t, allOutputsJson, "father.output")
		assert.Contains(t, allOutputsJson, "chick_1.output")
		assert.Contains(t, allOutputsJson, "chick_2.output")
	}
}

// Test error handling with non-existent stack output
func TestTgOutputErrorHandling(t *testing.T) {
	t.Parallel()

	// Create a temporary copy of the stack fixture
	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-error-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Initialize and apply terragrunt using stack commands
	_, err = TgStackInitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		ExtraArgs:        []string{"apply", "-auto-approve"},
	}
	_, err = TgStackRunE(t, applyOptions)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			ExtraArgs:        []string{"destroy", "-auto-approve"},
		}
		_, _ = TgStackRunE(t, destroyOptions)
	}()

	// Test that non-existent stack output returns error
	_, err = TgOutputE(t, options, "non_existent_output")
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "output")
}


