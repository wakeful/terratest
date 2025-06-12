package terragrunt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
)

func runTerragruntStackCommandE(t testing.TestingT, opts *Options, additionalArgs ...string) (string, error) {
	return runTerragruntStackSubCommandE(t, opts, "run", additionalArgs...)
}

func terragruntStackCommandE(t testing.TestingT, opts *Options, additionalArgs ...string) (string, error) {
	return runTerragruntStackSubCommandE(t, opts, "", additionalArgs...)
}

func runTerragruntStackSubCommandE(t testing.TestingT, opts *Options, subCommand string, additionalArgs ...string) (string, error) {
	args := []string{"stack"}
	if subCommand != "" {
		args = append(args, subCommand)
	}

	// Check for experimental flag support
	cmd := shell.Command{Command: opts.TerraformBinary, Args: []string{"-experiment", "stack"}}
	if err := shell.RunCommandE(t, cmd); err == nil {
		args = prepend(args, "-experiment", "stack")
	}

	options, args := terraform.GetCommonOptions(&opts.Options, args...)
	args = append(args, prepend(additionalArgs, "--")...)

	cmdExec := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.TerraformBinary, args)

	return retry.DoWithRetryableErrorsE(t, description, options.RetryableTerraformErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		s, err := shell.RunCommandAndGetOutputE(t, cmdExec)
		if err != nil {
			return s, err
		}
		if err := hasWarning(opts, s); err != nil {
			return s, err
		}
		return s, nil
	})
}

func prepend(args []string, arg ...string) []string {
	return append(arg, args...)
}

func hasWarning(opts *Options, out string) error {
	for k, v := range opts.WarningsAsErrors {
		str := fmt.Sprintf("\nWarning: %s[^\n]*\n", k)
		re, err := regexp.Compile(str)
		if err != nil {
			return fmt.Errorf("cannot compile regex for warning detection: %w", err)
		}
		m := re.FindAllString(out, -1)
		if len(m) == 0 {
			continue
		}
		return fmt.Errorf("warning(s) were found: %s:\n%s", v, strings.Join(m, ""))
	}
	return nil
}

func generateCommand(options *terraform.Options, args ...string) shell.Command {
	cmd := shell.Command{
		Command:    options.TerraformBinary,
		Args:       args,
		WorkingDir: options.TerraformDir,
		Env:        options.EnvVars,
		Logger:     options.Logger,
	}
	return cmd
}
