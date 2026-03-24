package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvArg_FromOSEnv(t *testing.T) {
	LocalEnv = map[string]string{}
	defer func() { LocalEnv = nil }()
	t.Setenv("GMMIT_TEST_VAR_OS", "os_value")

	result := GetEnvArg("GMMIT_TEST_VAR_OS")
	assert.Equal(t, "os_value", result)
}

func TestGetEnvArg_DefaultValue(t *testing.T) {
	LocalEnv = map[string]string{}
	defer func() { LocalEnv = nil }()
	os.Unsetenv("GMMIT_TEST_VAR_DEFAULT") //nolint:errcheck

	result := GetEnvArg("GMMIT_TEST_VAR_DEFAULT", "my_default")
	assert.Equal(t, "my_default", result)
}

func TestGetEnvArg_FromLocalEnv(t *testing.T) {
	LocalEnv = map[string]string{"GMMIT_TEST_LOCAL": "local_value"}
	defer func() { LocalEnv = nil }()

	result := GetEnvArg("GMMIT_TEST_LOCAL")
	assert.Equal(t, "local_value", result)
}

func TestGetEnvArg_LocalEnvBeatsOSEnv(t *testing.T) {
	t.Setenv("GMMIT_TEST_PRIORITY", "os_value")
	LocalEnv = map[string]string{"GMMIT_TEST_PRIORITY": "local_value"}
	defer func() { LocalEnv = nil }()

	result := GetEnvArg("GMMIT_TEST_PRIORITY")
	assert.Equal(t, "local_value", result)
}

func TestLoadEnvironment_NoFile(t *testing.T) {
	oldEnvFile := envFile
	envFile = "/tmp/non_existent_gmenv_test_file_xyz"
	defer func() {
		envFile = oldEnvFile
		LocalEnv = nil
	}()

	exited := false
	OsExit = func(int) { exited = true }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { LoadEnvironment() })

	require.False(t, exited)
	assert.NotNil(t, LocalEnv)
	assert.Empty(t, LocalEnv)
}

func TestLoadEnvironment_ValidFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "gmenv_test_*")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	_, err = tmpFile.WriteString("GMMIT_TEST_KEY=loaded_value\n")
	require.NoError(t, err)
	_ = tmpFile.Close()

	oldEnvFile := envFile
	envFile = tmpFile.Name()
	defer func() {
		envFile = oldEnvFile
		LocalEnv = nil
	}()

	exited := false
	OsExit = func(int) { exited = true }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { LoadEnvironment() })

	require.False(t, exited)
	assert.Equal(t, "loaded_value", LocalEnv["GMMIT_TEST_KEY"])
}
