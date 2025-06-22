package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgStackGenerate calls terragrunt stack run and returns stdout/stderr
func TgStackRun(t testing.TestingT, options *Options) string {
	out, err := TgStackRunE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgStackGenerateE calls terragrunt stack generate and returns stdout/stderr
func TgStackRunE(t testing.TestingT, options *Options) (string, error) {
	return terragruntStackCommandE(t, options, runArgs(options)...)
}

func runArgs(options *Options) []string {
	args := []string{"run"}

	args = append(args, options.ExtraArgs.Plan...)
	args = append(args, options.ExtraArgs.Apply...)
	args = append(args, options.ExtraArgs.Destroy...)

	// Append no-color option if needed
	if options.NoColor {
		args = append(args, "-no-color")
	}

	return args
}
