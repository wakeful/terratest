package terragrunt

import (
	"fmt"
	"slices"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
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
	return runTerragruntStackRunCommandE(t, options, runStackArgs(options)...)
}

// runTerragruntStackRunCommandE executes a terragrunt stack run command
// This is the specific implementation for stack run operations
func runTerragruntStackRunCommandE(t testing.TestingT, opts *Options, additionalArgs ...string) (string, error) {
	// Build the base command arguments starting with "stack run"
	commandArgs := []string{"stack", "run"}

	// Apply common terragrunt options and get the final command arguments
	terragruntOptions, finalArgs := GetCommonOptions(opts, commandArgs...)

	// Append additional arguments with "--" separator
	finalArgs = append(finalArgs, slices.Insert(additionalArgs, 0, "--")...)

	// Generate the final shell command
	execCommand := generateCommand(terragruntOptions, finalArgs...)
	commandDescription := fmt.Sprintf("%s %v", terragruntOptions.TerragruntBinary, finalArgs)

	// Execute the command with retry logic and error handling
	return retry.DoWithRetryableErrorsE(
		t,
		commandDescription,
		terragruntOptions.RetryableTerraformErrors,
		terragruntOptions.MaxRetries,
		terragruntOptions.TimeBetweenRetries,
		func() (string, error) {
			output, err := shell.RunCommandAndGetOutputE(t, execCommand)
			if err != nil {
				return output, err
			}

			// Check for warnings that should be treated as errors
			if warningErr := hasWarning(opts, output); warningErr != nil {
				return output, warningErr
			}

			return output, nil
		},
	)
}

func runStackArgs(options *Options) []string {
	args := []string{}

	args = append(args, options.ExtraArgs.Plan...)
	args = append(args, options.ExtraArgs.Apply...)
	args = append(args, options.ExtraArgs.Destroy...)

	// Append no-color option if needed
	if options.NoColor {
		args = append(args, "-no-color")
	}

	return args
}
