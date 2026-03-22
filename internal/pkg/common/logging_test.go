package common

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r) //nolint:errcheck
	os.Stdout = old
	return buf.String()
}

func TestInfo_Output(t *testing.T) {
	out := captureStdout(func() { Info("test message") })
	assert.Contains(t, out, "[I]")
	assert.Contains(t, out, "test message")
}

func TestInfoH_Output(t *testing.T) {
	out := captureStdout(func() { InfoH("infoh message") })
	assert.Contains(t, out, "[I]")
	assert.Contains(t, out, "infoh message")
}

func TestWarning_Output(t *testing.T) {
	out := captureStdout(func() { Warning("warn message") })
	assert.Contains(t, out, "[W]")
	assert.Contains(t, out, "warn message")
}

func TestError_Output(t *testing.T) {
	out := captureStdout(func() { Error("error message") })
	assert.Contains(t, out, "[E]")
	assert.Contains(t, out, "error message")
}

func TestDebug_NoOutput_WhenOff(t *testing.T) {
	LocalEnv = map[string]string{}
	defer func() { LocalEnv = nil }()
	os.Unsetenv("GMMIT_DEBUG") //nolint:errcheck

	out := captureStdout(func() { Debug("debug message") })
	assert.NotContains(t, out, "debug message")
}

func TestDebug_Output_WhenOn(t *testing.T) {
	LocalEnv = map[string]string{}
	defer func() { LocalEnv = nil }()
	t.Setenv("GMMIT_DEBUG", "true")

	out := captureStdout(func() { Debug("debug message") })
	assert.Contains(t, out, "debug message")
}

func TestPrintStartLine_Output(t *testing.T) {
	out := captureStdout(func() { PrintStartLine() })
	assert.Contains(t, out, "╔══[")
}

func TestPrintFailLine_CallsExit(t *testing.T) {
	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { PrintFailLine() })
	assert.Equal(t, 1, code)
}

func TestPrintFinalLine_CallsExit(t *testing.T) {
	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { PrintFinalLine() })
	assert.Equal(t, 0, code)
}

func TestDeleteLastLine_Output(t *testing.T) {
	out := captureStdout(func() { DeleteLastLine() })
	assert.Contains(t, out, "\033[1A")
}

func TestQuestion_Output(t *testing.T) {
	out := captureStdout(func() { Question("confirm?") })
	assert.Contains(t, out, "confirm?")
	assert.Contains(t, out, "Answer")
}
