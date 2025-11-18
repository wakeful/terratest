package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Init calls terragrunt init and return stdout/stderr
func Init(t testing.TestingT, options *Options) string {
	out, err := InitE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InitE calls terragrunt init and return stdout/stderr
func InitE(t testing.TestingT, options *Options) (string, error) {
	// Use regular terragrunt init command (not terragrunt stack init)
	return runTerragruntCommandE(t, options, "init", initArgs(options)...)
}

// initArgs builds the argument list for terragrunt init command.
// This function handles complex configuration that requires special formatting.
func initArgs(options *Options) []string {
	var args []string

	// Add complex configuration that requires special formatting
	// These are terraform-specific arguments that need special formatting
	args = append(args, terraform.FormatTerraformBackendConfigAsArgs(options.BackendConfig)...)
	args = append(args, terraform.FormatTerraformPluginDirAsArgs(options.PluginDir)...)
	return args
}
