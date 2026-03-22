package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRepoName_SSHGithub(t *testing.T) {
	name, provider, err := parseRepoName("git@github.com:owner/repo.git")
	require.NoError(t, err)
	assert.Equal(t, "owner/repo", name)
	assert.Equal(t, GIT_PROVIDER_GITHUB, provider)
}

func TestParseRepoName_SSHBitbucket(t *testing.T) {
	name, provider, err := parseRepoName("git@bitbucket.org:owner/repo.git")
	require.NoError(t, err)
	assert.Equal(t, "owner/repo", name)
	assert.Equal(t, GIT_PROVIDER_BITBUCKET, provider)
}

func TestParseRepoName_SSHGitlab(t *testing.T) {
	name, provider, err := parseRepoName("git@gitlab.com:owner/repo.git")
	require.NoError(t, err)
	assert.Equal(t, "owner/repo", name)
	assert.Equal(t, GIT_PROVIDER_GENERIC, provider)
}

func TestParseRepoName_GenericHost(t *testing.T) {
	name, provider, err := parseRepoName("git@custom-host.com:owner/myrepo.git")
	require.NoError(t, err)
	assert.Equal(t, "owner/myrepo", name)
	assert.Equal(t, GIT_PROVIDER_GENERIC, provider)
}

func TestParseRepoName_StripsDotGit(t *testing.T) {
	name, _, err := parseRepoName("git@github.com:owner/repo.git")
	require.NoError(t, err)
	assert.NotContains(t, name, ".git")
}

func TestParseRepoName_NoGitSuffix(t *testing.T) {
	name, provider, err := parseRepoName("git@github.com:owner/repo")
	require.NoError(t, err)
	assert.Equal(t, "owner/repo", name)
	assert.Equal(t, GIT_PROVIDER_GITHUB, provider)
}

func TestParseRepoName_NoColon_ReturnsError(t *testing.T) {
	_, _, err := parseRepoName("invalid-repo-string")
	assert.Error(t, err)
}

func TestParseRepoName_EmptyString_ReturnsError(t *testing.T) {
	_, _, err := parseRepoName("")
	assert.Error(t, err)
}
