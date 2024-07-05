package main

import (
	"fmt"
	"os"
	"os/exec"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

func main() {
	CheckArgs("<directory>")
	directory := os.Args[1]

	Info(directory)

	out, err := exec.Command("git","diff","--staged").Output()
    CheckIfError(err)
    
    fmt.Printf("git diff: %s\n", out)
}