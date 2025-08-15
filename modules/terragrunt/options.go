package terragrunt

import (
	"io"
	"os"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// Key concepts:
// - Options: Configure HOW the test framework executes tg (directories, retry logic, logging)
// - TerragruntArgs: Arguments for tg itself (e.g., --no-color for tg output)
// - TerraformArgs: Arguments passed to underlying terraform commands after -- separator
// - Use Options.TerragruntDir to specify WHERE to run tg
//
// Example:
//
//	// For init with terraform-specific flags
//	TgStackInitE(t, &Options{
//	    TerragruntDir: "/path/to/config",
//	    TerragruntArgs: []string{"--no-color"},
//	    TerraformArgs: []string{"-upgrade=true"},
//	})
//
//	// For stack run with terraform plan
//	TgStackRunE(t, &Options{
//	    TerragruntDir: "/path/to/config",
//	    TerragruntArgs: []string{"--no-color"},
//	    TerraformArgs: []string{"plan", "-out=tfplan"},
//	})
//
// Constants for test framework configuration and environment variables
const (
	DefaultTerragruntBinary = "terragrunt"
	NonInteractiveFlag      = "--non-interactive"
	TerragruntLogFormatKey  = "TG_LOG_FORMAT"
	TerragruntLogCustomKey  = "TG_LOG_CUSTOM_FORMAT"
	DefaultLogFormat        = "key-value"
	DefaultLogCustomFormat  = "%msg(color=disable)"
	ArgSeparator            = "--"
)

// Options represent the configuration options for tg test execution.
//
// This struct is divided into two clear categories:
//
// 1. TEST FRAMEWORK CONFIGURATION:
//   - Controls HOW the test framework executes tg
//   - Includes: binary paths, directories, retry logic, logging, environment
//   - These are NOT passed as command-line arguments to tg
//
// 2. TG COMMAND ARGUMENTS:
//   - All actual tg command-line arguments go in ExtraArgs []string
//   - This includes flags like -no-color, -upgrade, -reconfigure, etc.
//   - These ARE passed directly to the specific tg command being executed
//
// This separation eliminates confusion about which settings control the test
// framework vs which become tg command-line arguments.
type Options struct {
	// Test framework configuration (NOT passed to tg command line)
	TerragruntBinary string            // The tg binary to use (should be "terragrunt")
	TerragruntDir    string            // The directory containing the tg configuration
	EnvVars          map[string]string // Environment variables for command execution
	Logger           *logger.Logger    // Logger for command output

	// Test framework retry and error handling (NOT passed to tg command line)
	MaxRetries               int               // Maximum number of retries
	TimeBetweenRetries       time.Duration     // Time between retries
	RetryableTerraformErrors map[string]string // Retryable error patterns
	WarningsAsErrors         map[string]string // Warnings to treat as errors

	// Complex configuration that requires special formatting (NOT raw command-line args)
	BackendConfig map[string]interface{} // Backend configuration (formatted specially)
	PluginDir     string                 // Plugin directory (formatted specially)

	// Tg-specific command-line arguments (e.g., --no-color for tg itself)
	TerragruntArgs []string

	// Terraform command-line arguments to be passed after -- separator
	// These are passed directly to the underlying terraform commands
	TerraformArgs []string

	// Optional stdin to pass to Terraform commands
	Stdin io.Reader
}

// GetCommonOptions extracts common tg options and prepares arguments
// This is the tg-specific version of terraform.GetCommonOptions
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	// Set default binary if not specified
	if options.TerragruntBinary == "" {
		options.TerragruntBinary = DefaultTerragruntBinary
	}

	// Add tg-specific flags
	args = append(args, NonInteractiveFlag)

	// Set tg log formatting if not already set
	setTerragruntLogFormatting(options)

	return options, args
}

// GetArgsForCommand returns the appropriate arguments based on the command type
// It handles the separation of tg and terraform arguments
func GetArgsForCommand(options *Options, useArgSeparator bool) []string {
	var args []string

	// First add tg-specific arguments
	args = append(args, options.TerragruntArgs...)

	// Then add terraform arguments with separator if needed
	if len(options.TerraformArgs) > 0 {
		if useArgSeparator {
			args = append(args, ArgSeparator)
		}
		args = append(args, options.TerraformArgs...)
	}

	return args
}

// setTerragruntLogFormatting sets default log formatting for tg
// if it is not already set in options.EnvVars or OS environment vars
func setTerragruntLogFormatting(options *Options) {
	if options.EnvVars == nil {
		options.EnvVars = make(map[string]string)
	}

	_, inOpts := options.EnvVars[TerragruntLogFormatKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(TerragruntLogFormatKey)
		if !inEnv {
			// key-value format for tg logs to avoid colors and have plain form
			// https://terragrunt.gruntwork.io/docs/reference/cli-options/#terragrunt-log-format
			options.EnvVars[TerragruntLogFormatKey] = DefaultLogFormat
		}
	}

	_, inOpts = options.EnvVars[TerragruntLogCustomKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(TerragruntLogCustomKey)
		if !inEnv {
			options.EnvVars[TerragruntLogCustomKey] = DefaultLogCustomFormat
		}
	}
}
