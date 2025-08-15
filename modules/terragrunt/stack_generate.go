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
func generateStackArgs(options *Options) []string {
	args := []string{"generate"}

	// User-specified arguments are now handled by GetArgsForCommand
	// which properly separates TerragruntArgs and TerraformArgs

	return args
}
