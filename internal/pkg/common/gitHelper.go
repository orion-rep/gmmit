package common

import (
    "fmt"
	"strings"
	"errors"
)

// getRepoRawName returns the raw name of a repository from its full name
func getRepoRawName(repository string) (string, error) {
    repoGitName := strings.Split(repository, ":")
    if len(repoGitName) < 2 {
        return "", errors.New("Invalid repository name. " + repository + " could not be parsed properly")
    }
	// Removing '.git' from name
	repoName := strings.ReplaceAll(repoGitName[1], ".git", "")
    return repoName, nil
}

func GetDefaultBranch() (string){
    defaultBranch := strings.ReplaceAll(string(RunCommand("git", "rev-parse", "--abbrev-ref", "origin/HEAD")), "\n", "")
	Debug("Default branch: %s", defaultBranch)
    return defaultBranch
}

func GetCurrentBranch() (string){
    currentBranch := strings.ReplaceAll(string(RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")), "\n", "")
	Debug("Current branch: %s", currentBranch)
    return currentBranch
}

func GetRepositoryName() (string){
    remoteRepository := strings.ReplaceAll(string(RunCommand("git", "config", "--get", "remote.origin.url")), "\n", "")
	repositoryName, err := getRepoRawName(remoteRepository)
    CheckIfError(err)
	Debug("Repository name: %s", repositoryName)
    return repositoryName
}

func CalculateDiffBetweenBranches(branch1, branch2 string) (string){
	Debug("Checking changes between branch '%s' and '%s'", branch1, branch2)
    diff := string(RunCommand("git","diff",fmt.Sprintf("%s...%s", branch1, branch2)))
	return diff
}