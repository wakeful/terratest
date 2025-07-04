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

func generateStackArgs(options *Options) []string {
	args := []string{"generate"}

	// Append no-color option if needed
	if options.NoColor {
		args = append(args, "-no-color")
	}

	// Use Apply extra args for generate command as it's a similar operation
	if len(options.ExtraArgs.Apply) > 0 {
		args = append(args, options.ExtraArgs.Apply...)
	}
	return args
}
