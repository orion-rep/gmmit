package main

import (
	"os"
	"strings"
	
	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

func GetCommitContext()(string, string) {

	Info("Getting context information.")
	
	diff := string(RunCommand("git","diff","--staged"))

	if len(diff) <= 0 {
		Warning("Git diff returned no files.")
		Warning("Add some files to the staging area and run this command again.")
		os.Exit(0)
	}
	
	branch := strings.ReplaceAll(string(RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")), "\n", "")

	return diff, branch
    
}