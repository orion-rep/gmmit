package common

import (
	"strings"
	"errors"
)

// GetRepoRawName returns the raw name of a repository from its full name
func GetRepoRawName(repository string) (string, error) {
    repoGitName := strings.Split(repository, ":")
    if len(repoGitName) < 2 {
        return "", errors.New("Invalid repository name. " + repository + " could not be parsed properly")
    }
	// Removing '.git' from name
	repoName := strings.ReplaceAll(repoGitName[1], ".git", "")
    return repoName, nil
}