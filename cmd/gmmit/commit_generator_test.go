package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/assert"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
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

func makeResponse(text string) *genai.GenerateContentResponse {
	return &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{genai.Text(text)},
				},
			},
		},
	}
}

func TestGeneratePrompt_ContainsAllInputs(t *testing.T) {
	result := generatePrompt("feat: <desc>", "feature/123-my-feature", "diff content here")
	assert.Contains(t, result, "feat: <desc>")
	assert.Contains(t, result, "feature/123-my-feature")
	assert.Contains(t, result, "diff content here")
}

func TestGeneratePrompt_MentionsConventionalCommits(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff")
	assert.Contains(t, result, "Conventional Commits")
}

func TestGeneratePrompt_MentionsTicketID(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff")
	assert.Contains(t, result, "Ticket ID")
}

func TestGeneratePrompt_MentionsBranchName(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff")
	assert.Contains(t, result, "branch name")
}

func TestPrintHeader_ContainsVersion(t *testing.T) {
	Version = "v1.2.3"
	out := captureStdout(PrintHeader)
	assert.Contains(t, out, "v1.2.3")
}

func TestGetCommitContext_Success(t *testing.T) {
	callCount := 0
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		if callCount == 1 {
			return exec.Command("echo", "some staged diff")
		}
		return exec.Command("echo", "feature/123-my-branch")
	}
	defer func() { ExecCommand = exec.Command }()

	diff, branch := GetCommitContext()
	assert.Contains(t, diff, "some staged diff")
	assert.Equal(t, "feature/123-my-branch", branch)
}

func TestGetCommitContext_EmptyDiff_CallsExit(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true") // outputs nothing
	}
	defer func() { ExecCommand = exec.Command }()

	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	captureStdout(func() { GetCommitContext() })
	assert.Equal(t, 0, code) // PrintFinalLine exits with 0
}

func TestCreateCommit_NoVerify_False(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "1 file changed")
	}
	defer func() { ExecCommand = exec.Command }()

	f := false
	noVerifyFlag = &f
	captureStdout(func() { CreateCommit("feat: test commit") })
}

func TestCreateCommit_NoVerify_True(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "1 file changed")
	}
	defer func() { ExecCommand = exec.Command }()

	f := true
	noVerifyFlag = &f
	defer func() {
		ff := false
		noVerifyFlag = &ff
	}()
	captureStdout(func() { CreateCommit("feat: test commit --no-verify") })
}

func TestGenerateCommitMessage_ConfirmYes(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "1 file changed")
	}
	defer func() { ExecCommand = exec.Command }()

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse("feat: test commit message")
	}
	defer func() { runPromptFn = RunPrompt }()

	StdinReader = strings.NewReader("y\n")
	defer func() { StdinReader = os.Stdin }()

	f := false
	noVerifyFlag = &f
	commitStandard = "feat: <desc>"
	gitDiff = "some diff"
	gitBranch = "feature/123"

	captureStdout(func() { GenerateCommitMessage() })
}

func TestGenerateCommitMessage_Cancel(t *testing.T) {
	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse("feat: test commit message")
	}
	defer func() { runPromptFn = RunPrompt }()

	StdinReader = strings.NewReader("N\n")
	defer func() { StdinReader = os.Stdin }()

	var code int
	OsExit = func(c int) { code = c }
	defer func() { OsExit = os.Exit }()

	commitStandard = "feat: <desc>"
	gitDiff = "some diff"
	gitBranch = "feature/123"

	captureStdout(func() { GenerateCommitMessage() })
	assert.Equal(t, 0, code)
}

func TestGenerateCommitMessage_WithPush(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "ok")
	}
	defer func() { ExecCommand = exec.Command }()

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse("feat: test commit message")
	}
	defer func() { runPromptFn = RunPrompt }()

	StdinReader = strings.NewReader("y\n")
	defer func() { StdinReader = os.Stdin }()

	f := false
	noVerifyFlag = &f
	pt := true
	runCommitPush = &pt
	defer func() {
		pf := false
		runCommitPush = &pf
	}()
	commitStandard = "feat: <desc>"
	gitDiff = "some diff"
	gitBranch = "feature/123"

	captureStdout(func() { GenerateCommitMessage() })
}

func TestRunCommitGeneration(t *testing.T) {
	callCount := 0
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		switch callCount {
		case 1:
			return exec.Command("echo", "staged diff content")
		case 2:
			return exec.Command("echo", "feature/123")
		default:
			return exec.Command("echo", "1 file changed")
		}
	}
	defer func() { ExecCommand = exec.Command }()

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse("feat: test commit message")
	}
	defer func() { runPromptFn = RunPrompt }()

	StdinReader = strings.NewReader("y\n")
	defer func() { StdinReader = os.Stdin }()

	defer func() { LocalEnv = nil }()

	f := false
	noVerifyFlag = &f

	captureStdout(RunCommitGeneration)
}
