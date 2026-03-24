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

	gemini "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
	"gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	_ = w.Close()
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
	result := generatePrompt("feat: <desc>", "feature/123-my-feature", "diff content here", "", "")
	assert.Contains(t, result, "feat: <desc>")
	assert.Contains(t, result, "feature/123-my-feature")
	assert.Contains(t, result, "diff content here")
}

func TestGeneratePrompt_MentionsConventionalCommits(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff", "", "")
	assert.Contains(t, result, "Conventional Commits")
}

func TestGeneratePrompt_MentionsTicketID(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff", "", "")
	assert.Contains(t, result, "Ticket ID")
}

func TestGeneratePrompt_MentionsBranchName(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff", "", "")
	assert.Contains(t, result, "branch name")
}

func TestGeneratePrompt_WithTypeHint(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff", "fix", "")
	assert.Contains(t, result, `MUST be: "fix"`)
}

func TestGeneratePrompt_NoTypeHint(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff", "", "")
	assert.NotContains(t, result, "MUST be:")
}

func TestGeneratePrompt_WithHint(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff", "", "database migration")
	assert.Contains(t, result, `Focus on the following when writing the commit message: "database migration"`)
}

func TestGeneratePrompt_NoHint(t *testing.T) {
	result := generatePrompt("standard", "branch", "diff", "", "")
	assert.NotContains(t, result, "Focus on the following")
}

func TestValidateCommitType_Valid(t *testing.T) {
	for _, ct := range validCommitTypes {
		assert.True(t, validateCommitType(ct), "expected %q to be valid", ct)
	}
}

func TestValidateCommitType_Invalid(t *testing.T) {
	assert.False(t, validateCommitType("bad"))
	assert.False(t, validateCommitType(""))
	assert.False(t, validateCommitType("Fix"))
}

func TestRunCommitGeneration_InvalidType_CallsExit(t *testing.T) {
	bad := "bad"
	commitType = &bad
	defer func() { empty := ""; commitType = &empty }()

	var code int
	common.OsExit = func(c int) { code = c }
	defer func() { common.OsExit = os.Exit }()

	captureStdout(RunCommitGeneration)
	assert.Equal(t, 1, code)
}

func TestPrintHeader_ContainsVersion(t *testing.T) {
	Version = "v1.2.3"
	out := captureStdout(PrintHeader)
	assert.Contains(t, out, "v1.2.3")
}

func TestGetCommitContext_Success(t *testing.T) {
	callCount := 0
	common.ExecCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		if callCount == 1 {
			return exec.Command("echo", "some staged diff")
		}
		return exec.Command("echo", "feature/123-my-branch")
	}
	defer func() { common.ExecCommand = exec.Command }()

	diff, branch := GetCommitContext()
	assert.Contains(t, diff, "some staged diff")
	assert.Equal(t, "feature/123-my-branch", branch)
}

func TestGetCommitContext_EmptyDiff_CallsExit(t *testing.T) {
	common.ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true") // outputs nothing
	}
	defer func() { common.ExecCommand = exec.Command }()

	var code int
	common.OsExit = func(c int) { code = c }
	defer func() { common.OsExit = os.Exit }()

	captureStdout(func() { GetCommitContext() })
	assert.Equal(t, 0, code) // PrintFinalLine exits with 0
}

func TestCreateCommit_NoVerify_False(t *testing.T) {
	common.ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "1 file changed")
	}
	defer func() { common.ExecCommand = exec.Command }()

	f := false
	noVerifyFlag = &f
	captureStdout(func() { CreateCommit("feat: test commit") })
}

func TestCreateCommit_NoVerify_True(t *testing.T) {
	common.ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "1 file changed")
	}
	defer func() { common.ExecCommand = exec.Command }()

	f := true
	noVerifyFlag = &f
	defer func() {
		ff := false
		noVerifyFlag = &ff
	}()
	captureStdout(func() { CreateCommit("feat: test commit --no-verify") })
}

func TestGenerateCommitMessage_ConfirmYes(t *testing.T) {
	common.ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "1 file changed")
	}
	defer func() { common.ExecCommand = exec.Command }()

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse("feat: test commit message")
	}
	defer func() { runPromptFn = gemini.RunPrompt }()

	common.StdinReader = strings.NewReader("y\n")
	defer func() { common.StdinReader = os.Stdin }()

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
	defer func() { runPromptFn = gemini.RunPrompt }()

	common.StdinReader = strings.NewReader("N\n")
	defer func() { common.StdinReader = os.Stdin }()

	var code int
	common.OsExit = func(c int) { code = c }
	defer func() { common.OsExit = os.Exit }()

	commitStandard = "feat: <desc>"
	gitDiff = "some diff"
	gitBranch = "feature/123"

	captureStdout(func() { GenerateCommitMessage() })
	assert.Equal(t, 0, code)
}

func TestGenerateCommitMessage_WithPush(t *testing.T) {
	common.ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "ok")
	}
	defer func() { common.ExecCommand = exec.Command }()

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse("feat: test commit message")
	}
	defer func() { runPromptFn = gemini.RunPrompt }()

	common.StdinReader = strings.NewReader("y\n")
	defer func() { common.StdinReader = os.Stdin }()

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

func testRunCommitGenerationAddAll(t *testing.T, setFlag func(), resetFlag func()) {
	t.Helper()
	var capturedArgs [][]string
	callCount := 0
	common.ExecCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		capturedArgs = append(capturedArgs, append([]string{name}, args...))
		switch callCount {
		case 1: // git add .
			return exec.Command("true")
		case 2: // git diff --staged
			return exec.Command("echo", "staged diff content")
		case 3: // git rev-parse
			return exec.Command("echo", "feature/123")
		default: // git commit
			return exec.Command("echo", "1 file changed")
		}
	}
	defer func() { common.ExecCommand = exec.Command }()

	setFlag()
	defer resetFlag()

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse("feat: test commit message")
	}
	defer func() { runPromptFn = gemini.RunPrompt }()

	common.StdinReader = strings.NewReader("y\n")
	defer func() { common.StdinReader = os.Stdin }()

	defer func() { common.LocalEnv = nil }()

	f := false
	noVerifyFlag = &f

	captureStdout(RunCommitGeneration)

	assert.GreaterOrEqual(t, len(capturedArgs), 1)
	assert.Equal(t, []string{"git", "add", "."}, capturedArgs[0])
}

func TestRunCommitGeneration_AddAll(t *testing.T) {
	tr := true
	testRunCommitGenerationAddAll(t,
		func() { addAllFlag = &tr },
		func() { f := false; addAllFlag = &f },
	)
}

func TestRunCommitGeneration_AddAllShort(t *testing.T) {
	tr := true
	testRunCommitGenerationAddAll(t,
		func() { addAllShortFlag = &tr },
		func() { f := false; addAllShortFlag = &f },
	)
}

func TestRunCommitGeneration(t *testing.T) {
	callCount := 0
	common.ExecCommand = func(name string, args ...string) *exec.Cmd {
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
	defer func() { common.ExecCommand = exec.Command }()

	runPromptFn = func(p string) *genai.GenerateContentResponse {
		return makeResponse("feat: test commit message")
	}
	defer func() { runPromptFn = gemini.RunPrompt }()

	common.StdinReader = strings.NewReader("y\n")
	defer func() { common.StdinReader = os.Stdin }()

	defer func() { common.LocalEnv = nil }()

	f := false
	noVerifyFlag = &f

	captureStdout(RunCommitGeneration)
}
