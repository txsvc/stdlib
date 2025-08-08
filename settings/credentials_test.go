package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/stdlib/v2"
)

func TestCloneCredentials(t *testing.T) {
	cred := Credentials{
		ProjectID:    "p",
		ClientID:     "c",
		ClientSecret: "s",
		Token:        "t",
		Expires:      10,
	}
	dup := cred.Clone()
	assert.Equal(t, &cred, dup)
}

func TestValidation(t *testing.T) {
	cred1 := Credentials{}
	assert.False(t, cred1.IsValid())

	cred2 := Credentials{
		ProjectID:    "p",
		ClientID:     "c",
		ClientSecret: "s",
		Token:        "t",
		Expires:      10, // forces a fail ...
	}
	assert.False(t, cred2.IsValid())

	cred2.Expires = 0
	assert.True(t, cred2.IsValid())

	cred2.Token = ""
	assert.True(t, cred2.IsValid())

	cred2.ClientSecret = ""
	assert.True(t, cred2.IsValid())

	cred2.ClientID = ""
	assert.False(t, cred2.IsValid())

	cred2.Token = "t"
	assert.False(t, cred2.IsValid())

	cred2.Token = ""
	cred2.ClientSecret = "s"
	assert.False(t, cred2.IsValid())
}

func TestExpiration(t *testing.T) {
	cred := Credentials{
		ProjectID:    "p",
		ClientID:     "c",
		ClientSecret: "s",
		Token:        "t",
		Expires:      10,
	}
	assert.True(t, cred.Expired())

	cred.Expires = 0
	assert.False(t, cred.Expired())
}

func TestCredentialsFromEnvWithClientCredentials(t *testing.T) {
	expectedProjectID := "test-project-123"
	expectedClientID := "test-client-456"
	expectedClientSecret := "test-secret-789"

	_ = os.Setenv(ProjectID, expectedProjectID)
	_ = os.Setenv(ClientID, expectedClientID)
	_ = os.Setenv(ClientSecret, expectedClientSecret)
	defer func() {
		_ = os.Unsetenv(ProjectID)
		_ = os.Unsetenv(ClientID)
		_ = os.Unsetenv(ClientSecret)
	}()

	// Validate environment variables are set correctly
	assert.Equal(t, expectedProjectID, stdlib.GetString(ProjectID, ""))
	assert.Equal(t, expectedClientID, stdlib.GetString(ClientID, ""))
	assert.Equal(t, expectedClientSecret, stdlib.GetString(ClientSecret, ""))

	// Validate credentials object contains the correct values
	cred := CredentialsFromEnv()
	assert.NotNil(t, cred)
	assert.Equal(t, expectedProjectID, cred.ProjectID)
	assert.Equal(t, expectedClientID, cred.ClientID)
	assert.Equal(t, expectedClientSecret, cred.ClientSecret)
	assert.Empty(t, cred.Token) // Should be empty when using client credentials
}

func TestCredentialsFromEnvWithAccessToken(t *testing.T) {
	expectedProjectID := "test-project-123"
	expectedClientID := "test-client-456"
	expectedAccessToken := "test-token-xyz"

	_ = os.Setenv(ProjectID, expectedProjectID)
	_ = os.Setenv(ClientID, expectedClientID)
	_ = os.Setenv(AccessToken, expectedAccessToken)
	defer func() {
		_ = os.Unsetenv(ProjectID)
		_ = os.Unsetenv(ClientID)
		_ = os.Unsetenv(AccessToken)
	}()

	// Validate environment variables are set correctly
	assert.Equal(t, expectedProjectID, stdlib.GetString(ProjectID, ""))
	assert.Equal(t, expectedClientID, stdlib.GetString(ClientID, ""))
	assert.Equal(t, expectedAccessToken, stdlib.GetString(AccessToken, ""))

	// Validate credentials object contains the correct values
	cred := CredentialsFromEnv()
	assert.NotNil(t, cred)
	assert.Equal(t, expectedProjectID, cred.ProjectID)
	assert.Equal(t, expectedClientID, cred.ClientID)
	assert.Equal(t, expectedAccessToken, cred.Token)
	assert.Empty(t, cred.ClientSecret) // Should be empty when using access token
}

func TestCredentialsFromEnvEmpty(t *testing.T) {
	// Ensure no relevant env vars are set
	_ = os.Unsetenv(ProjectID)
	_ = os.Unsetenv(ClientID)
	_ = os.Unsetenv(ClientSecret)
	_ = os.Unsetenv(AccessToken)

	cred := CredentialsFromEnv()
	assert.NotNil(t, cred)
	assert.Empty(t, cred.ProjectID)
	assert.Empty(t, cred.ClientID)
	assert.Empty(t, cred.ClientSecret)
	assert.Empty(t, cred.Token)
}
