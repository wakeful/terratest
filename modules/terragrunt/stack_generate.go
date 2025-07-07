package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgStackGenerate calls terragrunt stack generate and returns stdout/stderr
func TgStackGenerate(t testing.TestingT, options *Options) string {
	out, err := TgStackGenerateE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgStackGenerateE calls terragrunt stack generate and returns stdout/stderr
func TgStackGenerateE(t testing.TestingT, options *Options) (string, error) {
	return terragruntStackCommandE(t, options, generateStackArgs(options)...)
}

// generateStackArgs builds the argument list for terragrunt stack generate command.
// All terragrunt command-line flags are now passed via ExtraArgs.
func generateStackArgs(options *Options) []string {
	args := []string{"generate"}

	// Add all user-specified terragrunt command-line arguments
	// This includes flags like -no-color, etc.
	args = append(args, options.ExtraArgs...)

	return args
}
