package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
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
	if options.TerraformBinary != "terragrunt" {
		return "", terraform.TgInvalidBinary(options.TerraformBinary)
	}
	return terragruntStackCommandE(t, options, generateArgs(options)...)
}
