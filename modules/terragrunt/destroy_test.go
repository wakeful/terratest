package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestDestroyAllNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	out := ApplyAll(t, options)
	require.Contains(t, out, "Hello, World")

	// Test that destroy completes successfully
	destroyOut := DestroyAll(t, options)
	require.NotEmpty(t, destroyOut, "Destroy output should not be empty")
}
