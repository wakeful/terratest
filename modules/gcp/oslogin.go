package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/oslogin/v1"
)

// ImportSSHKey will import an SSH key to GCP under the provided user identity.
// The `user` parameter should be the email address of the user.
// The `key` parameter should be the public key of the SSH key being uploaded.
// This will fail the test if there is an error.
func ImportSSHKey(t testing.TestingT, user, key string) {
	require.NoErrorf(t, ImportSSHKeyE(t, user, key), "Could not add SSH Key to user %s", user)
}

// ImportSSHKeyE will import an SSH key to GCP under the provided user identity.
// The `user` parameter should be the email address of the user.
// The `key` parameter should be the public key of the SSH key being uploaded.
func ImportSSHKeyE(t testing.TestingT, user, key string) error {
	return importProjectSSHKeyE(t, user, key, nil)
}

// ImportProjectSSHKey will import an SSH key to GCP under the provided user identity.
// The `user` parameter should be the email address of the user.
// The `key` parameter should be the public key of the SSH key being uploaded.
// The `projectID` parameter should be the chosen project ID.
// This will fail the test if there is an error.
func ImportProjectSSHKey(t testing.TestingT, user, key, projectID string) {
	require.NoErrorf(t, ImportProjectSSHKeyE(t, user, key, projectID), "Could not add SSH Key to user %s", user)
}

// ImportProjectSSHKeyE will import an SSH key to GCP under the provided user identity.
// The `user` parameter should be the email address of the user.
// The `key` parameter should be the public key of the SSH key being uploaded.
// The `projectID` parameter should be the chosen project ID.
func ImportProjectSSHKeyE(t testing.TestingT, user, key, projectID string) error {
	return importProjectSSHKeyE(t, user, key, &projectID)
}

func importProjectSSHKeyE(t testing.TestingT, user, key string, projectID *string) error {
	logger.Default.Logf(t, "Importing SSH key for user %s", user)

	ctx := context.Background()
	service, err := NewOSLoginServiceE(t)
	if err != nil {
		return err
	}

	parent := fmt.Sprintf("users/%s", user)

	sshPublicKey := &oslogin.SshPublicKey{
		Key: key,
	}

	req := service.Users.ImportSshPublicKey(parent, sshPublicKey)
	if projectID != nil {
		req = req.ProjectId(*projectID)
	}
	_, err = req.Context(ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

// DeleteSSHKey will delete an SSH key attached to the provided user identity.
// The `user` parameter should be the email address of the user.
// The `key` parameter should be the public key of the SSH key that was uploaded.
// This will fail the test if there is an error.
func DeleteSSHKey(t testing.TestingT, user, key string) {
	require.NoErrorf(t, DeleteSSHKeyE(t, user, key), "Could not delete SSH Key for user %s", user)
}

// DeleteSSHKeyE will delete an SSH key attached to the provided user identity.
// The `user` parameter should be the email address of the user.
// The `key` parameter should be the public key of the SSH key that was uploaded.
func DeleteSSHKeyE(t testing.TestingT, user, key string) error {
	logger.Default.Logf(t, "Deleting SSH key for user %s", user)

	ctx := context.Background()
	service, err := NewOSLoginServiceE(t)
	if err != nil {
		return err
	}

	loginProfile := GetLoginProfile(t, user)

	for _, v := range loginProfile.SshPublicKeys {
		if key == v.Key {
			path := fmt.Sprintf("users/%s/sshPublicKeys/%s", user, v.Fingerprint)
			_, err = service.Users.SshPublicKeys.Delete(path).Context(ctx).Do()
			break
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// GetLoginProfile will retrieve the login profile for a user's Google identity. The login profile is a combination of OS Login + gcloud SSH keys and POSIX
// accounts the user will appear as. Generally, this will only be the OS Login key + account, but `gcloud compute ssh` could create temporary keys and profiles.
// The `user` parameter should be the email address of the user.
// This will fail the test if there is an error.
func GetLoginProfile(t testing.TestingT, user string) *oslogin.LoginProfile {
	profile, err := GetLoginProfileE(t, user)
	require.NoErrorf(t, err, "Could not get login profile for user %s", user)

	return profile
}

// GetLoginProfileE will retrieve the login profile for a user's Google identity. The login profile is a combination of OS Login + gcloud SSH keys and POSIX
// accounts the user will appear as. Generally, this will only be the OS Login key + account, but `gcloud compute ssh` could create temporary keys and profiles.
// The `user` parameter should be the email address of the user.
func GetLoginProfileE(t testing.TestingT, user string) (*oslogin.LoginProfile, error) {
	logger.Default.Logf(t, "Getting login profile for user %s", user)

	ctx := context.Background()
	service, err := NewOSLoginServiceE(t)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("users/%s", user)

	profile, err := service.Users.GetLoginProfile(name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// NewOSLoginServiceE creates a new OS Login service, which is used to make OS Login API calls.
func NewOSLoginServiceE(t testing.TestingT) (*oslogin.Service, error) {
	ctx := context.Background()

	if ts, ok := getStaticTokenSource(); ok {
		return oslogin.New(oauth2.NewClient(ctx, ts))
	}

	client, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("Failed to get default client: %v", err)
	}

	service, err := oslogin.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}
