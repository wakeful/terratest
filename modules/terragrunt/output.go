package terragrunt

import (
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgOutput calls terragrunt stack output for the given variable and returns its value as a string
func TgOutput(t testing.TestingT, options *Options, key string) string {
	out, err := TgOutputE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgOutputE calls terragrunt stack output for the given variable and returns its value as a string
func TgOutputE(t testing.TestingT, options *Options, key string) (string, error) {
	rawOutput, err := runTerragruntStackCommandE(t, options, "output", outputArgs(options, key)...)
	if err != nil {
		return "", err
	}
	return cleanTerragruntOutput(rawOutput), nil
}

// outputArgs builds the argument list for terragrunt stack output command
func outputArgs(options *Options, key string) []string {
	args := []string{"-no-color"}

	// Add all user-specified terragrunt command-line arguments first
	args = append(args, options.ExtraArgs...)

	// Add the key last, if provided
	if key != "" {
		args = append(args, key)
	}

	return args
}

// cleanTerragruntOutput extracts the actual output value from terragrunt stack's verbose output
func cleanTerragruntOutput(rawOutput string) string {
	lines := strings.Split(rawOutput, "\n")

	var outputLines []string
	inJsonBlock := false
	braceCount := 0

	// Process lines to extract the actual output
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Skip terragrunt log lines
		if strings.HasPrefix(trimmedLine, "time=") {
			continue
		}

		if trimmedLine == "" {
			continue
		}

		// Check if this line starts a JSON block
		if strings.HasPrefix(trimmedLine, "{") && !inJsonBlock {
			inJsonBlock = true
			braceCount = strings.Count(trimmedLine, "{") - strings.Count(trimmedLine, "}")
			outputLines = append(outputLines, trimmedLine)
		} else if inJsonBlock {
			// We're in a JSON block, count braces
			braceCount += strings.Count(trimmedLine, "{") - strings.Count(trimmedLine, "}")
			outputLines = append(outputLines, trimmedLine)

			// If braces are balanced, we're done with JSON
			if braceCount <= 0 {
				break
			}
		} else {
			// Not in JSON block, this is likely a simple value
			// Remove surrounding quotes if present
			if strings.HasPrefix(trimmedLine, "\"") && strings.HasSuffix(trimmedLine, "\"") {
				trimmedLine = strings.Trim(trimmedLine, "\"")
			}
			return trimmedLine
		}
	}

	// If we collected JSON lines, return them joined
	if len(outputLines) > 0 {
		return strings.Join(outputLines, "")
	}

	return rawOutput
}
