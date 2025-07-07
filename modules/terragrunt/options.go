package terragrunt

import (
	"os"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// Key concepts:
// - Options: Configure HOW the test framework executes terragrunt (directories, retry logic, logging)
// - ExtraArgs: Specify ALL command-line arguments passed to the specific terragrunt command
// - Use Options.TerragruntDir to specify WHERE to run terragrunt
// - Use ExtraArgs to pass ALL command-line arguments (including -no-color, -upgrade, etc.)
//
// Example:
//
//	// For init
//	TgStackInitE(t, &Options{
//	    TerragruntDir: "/path/to/config",
//	    ExtraArgs: []string{"-upgrade=true", "-no-color"},
//	})
//
//	// For generate
//	TgStackGenerateE(t, &Options{
//	    TerragruntDir: "/path/to/config",
//	    ExtraArgs: []string{"-no-color"},
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

// Options represent the configuration options for terragrunt test execution.
//
// This struct is divided into two clear categories:
//
// 1. TEST FRAMEWORK CONFIGURATION:
//   - Controls HOW the test framework executes terragrunt
//   - Includes: binary paths, directories, retry logic, logging, environment
//   - These are NOT passed as command-line arguments to terragrunt
//
// 2. TERRAGRUNT COMMAND ARGUMENTS:
//   - All actual terragrunt command-line arguments go in ExtraArgs []string
//   - This includes flags like -no-color, -upgrade, -reconfigure, etc.
//   - These ARE passed directly to the specific terragrunt command being executed
//
// This separation eliminates confusion about which settings control the test
// framework vs which become terragrunt command-line arguments.
type Options struct {
	// Test framework configuration (NOT passed to terragrunt command line)
	TerragruntBinary string            // The terragrunt binary to use (should be "terragrunt")
	TerragruntDir    string            // The directory containing the terragrunt configuration
	EnvVars          map[string]string // Environment variables for command execution
	Logger           *logger.Logger    // Logger for command output

	// Test framework retry and error handling (NOT passed to terragrunt command line)
	MaxRetries               int               // Maximum number of retries
	TimeBetweenRetries       time.Duration     // Time between retries
	RetryableTerraformErrors map[string]string // Retryable error patterns
	WarningsAsErrors         map[string]string // Warnings to treat as errors

	// Complex configuration that requires special formatting (NOT raw command-line args)
	BackendConfig map[string]interface{} // Backend configuration (formatted specially)
	PluginDir     string                 // Plugin directory (formatted specially)

	// All terragrunt command-line arguments for the specific command being executed
	ExtraArgs []string
}

// GetCommonOptions extracts common terragrunt options and prepares arguments
// This is the terragrunt-specific version of terraform.GetCommonOptions
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	// Set default binary if not specified
	if options.TerragruntBinary == "" {
		options.TerragruntBinary = DefaultTerragruntBinary
	}

	// Add terragrunt-specific flags
	args = append(args, NonInteractiveFlag)

	// Set terragrunt log formatting if not already set
	setTerragruntLogFormatting(options)

	return options, args
}

// setTerragruntLogFormatting sets default log formatting for terragrunt
// if it is not already set in options.EnvVars or OS environment vars
func setTerragruntLogFormatting(options *Options) {
	if options.EnvVars == nil {
		options.EnvVars = make(map[string]string)
	}

	_, inOpts := options.EnvVars[TerragruntLogFormatKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(TerragruntLogFormatKey)
		if !inEnv {
			// key-value format for terragrunt logs to avoid colors and have plain form
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
