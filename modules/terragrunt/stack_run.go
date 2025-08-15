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
// This function is now just returning an empty slice since arguments
// are handled by GetArgsForCommand in cmd.go
func runStackArgs(options *Options) []string {
	// Arguments are now handled by GetArgsForCommand which properly
	// separates TerragruntArgs and TerraformArgs
	return []string{}
}
