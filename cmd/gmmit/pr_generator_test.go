package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/assert"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

func TestCreatePROnGithub_Success(t *testing.T) {
	defer func() { LocalEnv = nil }()
	t.Setenv("GMMIT_GH_USER", "testuser")
	t.Setenv("GMMIT_GH_PASS", "testtoken")

	callPostFn = func(url string, payload interface{}, user string, pass string) ([]byte, int, error) {
		return []byte(`{"html_url":"https://github.com/owner/repo/pull/1"}`), 201, nil
	}
	defer func() { callPostFn = CallPost }()

	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { ExecCommand = exec.Command }()

	result := captureStdout(func() {
		url := createPROnGithub("feat: test", "description", "feature/123", "main", "owner/repo")
		assert.Equal(t, "https://github.com/owner/repo/pull/1", url)
	})
	assert.Contains(t, result, "PR URL")
}

func TestCreatePROnBitbucket_Success(t *testing.T) {
	defer func() { LocalEnv = nil }()
	t.Setenv("GMMIT_BB_USER", "bbuser")
	t.Setenv("GMMIT_BB_PASS", "bbpass")

	callPostFn = func(url string, payload interface{}, user string, pass string) ([]byte, int, error) {
		return []byte(`{"links":{"html":{"href":"https://bitbucket.org/owner/repo/pull-requests/1"}}}`), 201, nil
	}
	defer func() { callPostFn = CallPost }()

	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { ExecCommand = exec.Command }()

	result := captureStdout(func() {
		url := createPROnBitbucket("feat: test", "description", "feature/123", "owner/repo")
		assert.Equal(t, "https://bitbucket.org/owner/repo/pull-requests/1", url)
	})
	assert.Contains(t, result, "PR URL")
}

func TestGetPRContext_Success(t *testing.T) {
	callCount := 0
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		switch callCount {
		case 1:
			return exec.Command("echo", "git@github.com:owner/repo.git") // GetRepositoryData
		case 2:
			return exec.Command("echo", "origin/main") // GetDefaultBranch
		case 3:
			return exec.Command("echo", "feature/123") // GetCurrentBranch
		default:
			return exec.Command("echo", "some diff content") // CalculateDiffBetweenBranches
		}
	}
	defer func() { ExecCommand = exec.Command }()

	captureStdout(getPRContext)

	assert.Equal(t, "owner/repo", repositoryName)
	assert.Equal(t, "Github", repositoryProvider)
	assert.Equal(t, "main", gitDefaultBranch)
	assert.Equal(t, "feature/123", gitPRBranch)
	assert.Contains(t, gitPRDiff, "some diff")
}

func TestGeneratePRMessage_GenericProvider(t *testing.T) {
	repositoryProvider = "Generic"
	gitPRBranch = "feature/123"
	gitPRDiff = "some diff content"

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse(`{"title":"feat: test","description":"test description"}`)
	}
	defer func() { runPromptFn = RunPrompt }()

	// "N" cancels confirmCopyClipboard, which calls PrintFinalLine
	StdinReader = strings.NewReader("N\n")
	defer func() { StdinReader = os.Stdin }()

	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	defer func() { LocalEnv = nil }()

	captureStdout(generatePRMessage)
	assert.Equal(t, 0, code)
}

func TestGeneratePRMessage_GithubProvider_Cancel(t *testing.T) {
	repositoryProvider = "Github"
	gitPRBranch = "feature/123"
	gitPRDiff = "some diff content"
	repositoryName = "owner/repo"
	gitDefaultBranch = "main"

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse(`{"title":"feat: test","description":"test description"}`)
	}
	defer func() { runPromptFn = RunPrompt }()

	// "N" cancels PR creation, calls confirmCopyClipboard, then "N" again cancels that
	StdinReader = strings.NewReader("N\nN\n")
	defer func() { StdinReader = os.Stdin }()

	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	defer func() { LocalEnv = nil }()

	captureStdout(generatePRMessage)
	assert.Equal(t, 0, code)
}

func TestRunPRGeneration(t *testing.T) {
	callCount := 0
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		switch callCount {
		case 1:
			return exec.Command("echo", "git@github.com:owner/repo.git")
		case 2:
			return exec.Command("echo", "origin/main")
		case 3:
			return exec.Command("echo", "feature/123")
		default:
			return exec.Command("echo", "some diff content")
		}
	}
	defer func() { ExecCommand = exec.Command }()

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse(`{"title":"feat: test","description":"test description"}`)
	}
	defer func() { runPromptFn = RunPrompt }()

	StdinReader = strings.NewReader("N\nN\n")
	defer func() { StdinReader = os.Stdin }()

	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	captureStdout(RunPRGeneration)
	assert.Equal(t, 0, code)
}

func TestConfirmCopyClipboard_Cancel(t *testing.T) {
	StdinReader = strings.NewReader("N\n")
	defer func() { StdinReader = os.Stdin }()

	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { confirmCopyClipboard("test description") })
	assert.Equal(t, 0, code)
}
