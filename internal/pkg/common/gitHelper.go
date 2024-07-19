package common

import (
    "fmt"
	"strings"
	"errors"
)

func parseRepoName(repository string) (string, string, error) {
    repoGitName := strings.Split(repository, ":")
    if len(repoGitName) < 2 {
        return "", "", errors.New("Invalid repository name. " + repository + " could not be parsed properly")
    }
	// Removing '.git' from name
	repoName := strings.ReplaceAll(repoGitName[1], ".git", "")
    Debug("Repository name: %s", repoName)

    repoProvider := "Generic"
    if strings.Contains(repoGitName[0], "bitbucket") {
        repoProvider = "Bitbucket"
    } else if strings.Contains(repoGitName[0], "github") {
        // repoProvider = "Github" // Not supported yes
        Debug("Github provider is not supported yes")
    } else if strings.Contains(repoGitName[0], "gitlab") {
        // repoProvider = "Gitlab" // Not supported yet
        Debug("Gitlab provider is not supported yet")
    }
    Debug("Repository provider: %s", repoProvider)

    return repoName, repoProvider, nil
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

func GetRepositoryData() (string, string){
    remoteRepository := strings.ReplaceAll(string(RunCommand("git", "config", "--get", "remote.origin.url")), "\n", "")
	repositoryName, repositoryProvider, err := parseRepoName(remoteRepository)
    CheckIfError(err)
	return repositoryName, repositoryProvider
}

func CalculateDiffBetweenBranches(branch1, branch2 string) (string){
	Debug("Checking changes between branch '%s' and '%s'", branch1, branch2)
    diff := string(RunCommand("git","diff",fmt.Sprintf("%s...%s", branch1, branch2)))
	return diff
}