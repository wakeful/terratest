package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTerragruntInit(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	out, err := TgStackInitE(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		ExtraArgs:        []string{"-upgrade=true"}, // Common init flag
	})
	require.NoError(t, err)
	require.Contains(t, out, "Terraform has been successfully initialized!")
}

func TestTerragruntInitWithInvalidConfig(t *testing.T) {
	t.Parallel()
	// Test error handling when terragrunt.hcl has invalid HCL syntax
	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init-error", t.Name())
	require.NoError(t, err)

	// This should fail due to invalid HCL syntax in terragrunt.hcl
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		ExtraArgs:        []string{"-upgrade=true"}, // Common init flag
	})
	require.Error(t, err)
	// The error should contain information about the HCL parsing error
	require.Contains(t, err.Error(), "Missing expression")
}
