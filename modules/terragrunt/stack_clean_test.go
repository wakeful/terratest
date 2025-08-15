package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTgStackClean(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")

	// First generate the stack to create .terragrunt-stack directory
	_, err = TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Verify that the .terragrunt-stack directory was created
	require.DirExists(t, stackDir)

	// Clean the stack
	out, err := TgStackCleanE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Verify clean command produced expected output
	require.Contains(t, out, "Deleting stack directory")

	// Verify that the .terragrunt-stack directory was removed
	require.NoDirExists(t, stackDir)
}

func TestTgStackCleanNonExistentStack(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")

	// Verify that the .terragrunt-stack directory doesn't exist
	require.NoDirExists(t, stackDir)

	// Clean should succeed even if .terragrunt-stack doesn't exist
	_, err = TgStackCleanE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)
}

func TestTgStackCleanAfterRun(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")

	// Initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Run plan to generate the stack
	_, err = TgStackRunE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"plan"},
	})
	require.NoError(t, err)

	// Verify that the .terragrunt-stack directory was created
	require.DirExists(t, stackDir)

	// Clean the stack
	out, err := TgStackCleanE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Verify clean command produced expected output
	require.Contains(t, out, "Deleting stack directory")

	// Verify that the .terragrunt-stack directory was removed
	require.NoDirExists(t, stackDir)
}