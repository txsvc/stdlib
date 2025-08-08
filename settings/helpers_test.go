package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testCredentialFile = "test.json"

func TestWriteReadSettings(t *testing.T) {
	settings1 := &DialSettings{
		Endpoint: "x",
		//DefaultEndpoint: "X",
		Scopes:        []string{"a", "b"},
		DefaultScopes: []string{"A", "B"},
		UserAgent:     "agent",
	}
	settings1.SetOption("FOO", "x")
	settings1.SetOption("BAR", "x")

	err := WriteDialSettings(settings1, testCredentialFile)
	assert.NoError(t, err)

	settings2, err := ReadDialSettings(testCredentialFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, settings2)
	assert.Equal(t, settings1, settings2)

	// cleanup
	_ = os.Remove(testCredentialFile)
}

func TestWriteReadCredentials(t *testing.T) {
	cred1 := &Credentials{
		ProjectID: "project",
		ClientID:  "client",
		Token:     "token",
		Expires:   42,
	}

	err := WriteCredentials(cred1, testCredentialFile)
	assert.NoError(t, err)

	cred2, err := ReadCredentials(testCredentialFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, cred2)
	assert.Equal(t, cred1, cred2)

	// cleanup
	_ = os.Remove(testCredentialFile)
}

func TestReadNonexistentDialSettings(t *testing.T) {
	_, err := ReadDialSettings("nonexistent.json")
	assert.Error(t, err)
}

func TestReadNonexistentCredentials(t *testing.T) {
	_, err := ReadCredentials("nonexistent.json")
	assert.Error(t, err)
}

func TestReadInvalidJsonDialSettings(t *testing.T) {
	invalidFile := "invalid_dial.json"
	err := os.WriteFile(invalidFile, []byte("invalid json"), 0644)
	assert.NoError(t, err)
	defer func() { _ = os.Remove(invalidFile) }()

	_, err = ReadDialSettings(invalidFile)
	assert.Error(t, err)
}

func TestReadInvalidJsonCredentials(t *testing.T) {
	invalidFile := "invalid_cred.json"
	err := os.WriteFile(invalidFile, []byte("invalid json"), 0644)
	assert.NoError(t, err)
	defer func() { _ = os.Remove(invalidFile) }()

	_, err = ReadCredentials(invalidFile)
	assert.Error(t, err)
}

func TestWriteDialSettingsWithDirectoryCreation(t *testing.T) {
	nestedPath := "test_dir/subdir/settings.json"
	defer func() { _ = os.RemoveAll("test_dir") }()

	settings1 := &DialSettings{
		Endpoint:      "test-endpoint",
		Scopes:        []string{"scope1", "scope2"},
		DefaultScopes: []string{"default1"},
		UserAgent:     "test-agent",
	}

	err := WriteDialSettings(settings1, nestedPath)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(nestedPath)
	assert.NoError(t, err)

	// Verify content is correct
	settings2, err := ReadDialSettings(nestedPath)
	assert.NoError(t, err)
	assert.Equal(t, settings1, settings2)
}

func TestWriteCredentialsWithDirectoryCreation(t *testing.T) {
	nestedPath := "test_dir2/subdir/creds.json"
	defer func() { _ = os.RemoveAll("test_dir2") }()

	cred1 := &Credentials{
		ProjectID:    "test-project",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		Token:        "test-token",
		Expires:      1234567890,
	}

	err := WriteCredentials(cred1, nestedPath)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(nestedPath)
	assert.NoError(t, err)

	// Verify content is correct
	cred2, err := ReadCredentials(nestedPath)
	assert.NoError(t, err)
	assert.Equal(t, cred1, cred2)
}

func TestWriteReadEmptyDialSettings(t *testing.T) {
	emptySettings := &DialSettings{}

	err := WriteDialSettings(emptySettings, "empty_dial.json")
	assert.NoError(t, err)
	defer func() { _ = os.Remove("empty_dial.json") }()

	readSettings, err := ReadDialSettings("empty_dial.json")
	assert.NoError(t, err)
	assert.Equal(t, emptySettings, readSettings)
}

func TestWriteReadEmptyCredentials(t *testing.T) {
	emptyCreds := &Credentials{}

	err := WriteCredentials(emptyCreds, "empty_creds.json")
	assert.NoError(t, err)
	defer func() { _ = os.Remove("empty_creds.json") }()

	readCreds, err := ReadCredentials("empty_creds.json")
	assert.NoError(t, err)
	assert.Equal(t, emptyCreds, readCreds)
}

func TestDialSettingsWithOptions(t *testing.T) {
	settings1 := &DialSettings{
		Endpoint:  "test-endpoint",
		UserAgent: "test-agent",
	}
	settings1.SetOption("CustomOption1", "value1")
	settings1.SetOption("CustomOption2", "42")
	settings1.SetOption("CustomOption3", "true")

	err := WriteDialSettings(settings1, "settings_with_options.json")
	assert.NoError(t, err)
	defer func() { _ = os.Remove("settings_with_options.json") }()

	settings2, err := ReadDialSettings("settings_with_options.json")
	assert.NoError(t, err)
	assert.Equal(t, settings1, settings2)

	// Verify options are preserved
	assert.Equal(t, "value1", settings2.GetOption("CustomOption1"))
	assert.Equal(t, "42", settings2.GetOption("CustomOption2"))
	assert.Equal(t, "true", settings2.GetOption("CustomOption3"))
}
