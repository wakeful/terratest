package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestTerragruntStackGenerate(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		Options: terraform.Options{
			TerraformDir:    path.Join(testFolder, "live"),
			TerraformBinary: "terragrunt",
		},
	})
	require.NoError(t, err)

	// Then generate the stack
	_, err = TgStackGenerateE(t, &Options{
		Options: terraform.Options{
			TerraformDir:    path.Join(testFolder, "live"),
			TerraformBinary: "terragrunt",
		},
	})
	require.NoError(t, err)
}
