package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestTerragruntStackInit(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	out, err := TgStackInitE(t, &Options{
		Options: terraform.Options{
			TerraformDir:    path.Join(testFolder, "live"),
			TerraformBinary: "terragrunt",
		},
	})
	require.NoError(t, err)
	require.Contains(t, out, ".terragrunt-stack")
	require.Contains(t, out, "has been successfully initialized!")
}
