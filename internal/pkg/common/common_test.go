package common

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand_ReturnsStdout(t *testing.T) {
	output := RunCommand("echo", "hello")
	assert.Contains(t, output, "hello")
}

func TestAskConfirmation_Yes(t *testing.T) {
	StdinReader = strings.NewReader("y\n")
	defer func() { StdinReader = os.Stdin }()

	result := captureStdout(func() {
		// AskConfirmation writes to stdout (Question + DeleteLastLine); capture it to keep test output clean
	})
	_ = result

	StdinReader = strings.NewReader("y\n")
	var got int
	captureStdout(func() { got = AskConfirmation("test?") })
	assert.Equal(t, 1, got)
}

func TestAskConfirmation_YesFull(t *testing.T) {
	StdinReader = strings.NewReader("yes\n")
	defer func() { StdinReader = os.Stdin }()

	var got int
	captureStdout(func() { got = AskConfirmation("test?") })
	assert.Equal(t, 1, got)
}

func TestAskConfirmation_Retry(t *testing.T) {
	StdinReader = strings.NewReader("r\n")
	defer func() { StdinReader = os.Stdin }()

	var got int
	captureStdout(func() { got = AskConfirmation("test?") })
	assert.Equal(t, 2, got)
}

func TestAskConfirmation_No(t *testing.T) {
	StdinReader = strings.NewReader("N\n")
	defer func() { StdinReader = os.Stdin }()

	var got int
	captureStdout(func() { got = AskConfirmation("test?") })
	assert.Equal(t, 0, got)
}

func TestAskConfirmation_Default(t *testing.T) {
	StdinReader = strings.NewReader("anything\n")
	defer func() { StdinReader = os.Stdin }()

	var got int
	captureStdout(func() { got = AskConfirmation("test?") })
	assert.Equal(t, 0, got)
}

func TestAskConfirmation_AutoConfirm(t *testing.T) {
	AutoConfirm = true
	defer func() { AutoConfirm = false }()

	// StdinReader should never be read when AutoConfirm is set
	StdinReader = strings.NewReader("")

	var got int
	captureStdout(func() { got = AskConfirmation("test?") })
	assert.Equal(t, 1, got)
}

func TestCheckIfError_NilError(t *testing.T) {
	// Must not call OsExit for nil error
	exited := false
	OsExit = func(int) { exited = true }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { CheckIfError(nil) })
	assert.False(t, exited)
}

func TestCheckIfError_WithError(t *testing.T) {
	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { CheckIfError(errors.New("something broke")) })
	assert.Equal(t, 1, code)
}

func TestCheckIfError_WithContext(t *testing.T) {
	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { CheckIfError(errors.New("err"), "extra context") })
	assert.Equal(t, 1, code)
}

func TestCheckArgs_Sufficient(t *testing.T) {
	// In test context os.Args has at least the binary name; CheckArgs with zero required args must not exit.
	exited := false
	OsExit = func(int) { exited = true }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { CheckArgs() })
	assert.False(t, exited)
}

func TestCheckArgs_Insufficient(t *testing.T) {
	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	oldArgs := os.Args
	os.Args = []string{"gmmit"} // only 1 element
	defer func() { os.Args = oldArgs }()

	captureStdout(func() { CheckArgs("required-arg") })
	assert.Equal(t, 1, code)
}

func TestCommandExists_ExistingCommand(t *testing.T) {
	exited := false
	OsExit = func(int) { exited = true }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { CommandExists("ls") })
	assert.False(t, exited)
}

func TestOpenURL_Success(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { ExecCommand = exec.Command }()

	err := OpenURL("https://example.com")
	assert.NoError(t, err)
}
