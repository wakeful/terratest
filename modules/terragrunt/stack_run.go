package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgStackRun calls terragrunt stack run and returns stdout/stderr
func TgStackRun(t testing.TestingT, options *Options) string {
	out, err := TgStackRunE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgStackRunE calls terragrunt stack run and returns stdout/stderr
func TgStackRunE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "run", runStackArgs(options)...)
}

// runStackArgs builds the argument list for terragrunt stack run command.
// All terragrunt command-line flags are now passed via ExtraArgs.
func runStackArgs(options *Options) []string {
	// Return all user-specified terragrunt command-line arguments
	// The user passes the specific args they need for their stack run operation
	return options.ExtraArgs
}
