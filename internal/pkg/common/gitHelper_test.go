package common

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentBranch(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "main")
	}
	defer func() { ExecCommand = exec.Command }()

	result := GetCurrentBranch()
	assert.Equal(t, "main", result)
}

func TestGetDefaultBranch(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "origin/main")
	}
	defer func() { ExecCommand = exec.Command }()

	result := GetDefaultBranch()
	assert.Equal(t, "main", result)
}

func TestGetRepositoryData_Github(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "git@github.com:owner/repo.git")
	}
	defer func() { ExecCommand = exec.Command }()

	name, provider := GetRepositoryData()
	assert.Equal(t, "owner/repo", name)
	assert.Equal(t, GIT_PROVIDER_GITHUB, provider)
}

func TestGetRepositoryData_Bitbucket(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "git@bitbucket.org:owner/repo.git")
	}
	defer func() { ExecCommand = exec.Command }()

	name, provider := GetRepositoryData()
	assert.Equal(t, "owner/repo", name)
	assert.Equal(t, GIT_PROVIDER_BITBUCKET, provider)
}

func TestCalculateDiffBetweenBranches(t *testing.T) {
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", "diff output")
	}
	defer func() { ExecCommand = exec.Command }()

	result := CalculateDiffBetweenBranches("main", "feature/123")
	assert.Contains(t, result, "diff output")
}
