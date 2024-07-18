package common

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"log"
	"errors"

	"github.com/joho/godotenv"
)

var localEnv map[string]string

var envFile = string(os.Getenv("HOME")) + "/.gmenv"

func GetEnvArg(name string, defaultValue ...string) (string) {
	var value = ""
	if val, ok := localEnv[name]; ok {
		value = val
	} else if envVal, exists := os.LookupEnv(name); exists {
		value = envVal
	} else if len(defaultValue) > 0 {
		value = defaultValue[0]
	} else {
		switch option := AskConfirmation("Config var '" + name + "' is not define. Do you want to set it? [y/N]"); option {
		case 1:
			defineEnvArg(name)
			return GetEnvArg(name, defaultValue...)
		default:
			Error(fmt.Sprintf("Config var '%s' is not present", name))
			os.Exit(1)
		}
	}
	return value
}

func LoadEnvironment()() {
	Debug("Checking environment packages and tools")
	CommandExists("git")
	Debug("Loading environment variables from ~/.gmenv file")

	if _, err := os.Stat(envFile); errors.Is(err, os.ErrNotExist) {
		Warning("Env file %s doesn't exist", envFile)
		localEnv = make(map[string]string)
		return 
	}

	env, err := godotenv.Read(envFile)
	if err != nil {
		Error("Error loading ~/.gmenv file")
		Error(err.Error())
		os.Exit(1)
	}

	localEnv = env
}

func defineEnvArg(name string) {
	Question("Define '" + name + "' value:")
	reader := bufio.NewReader(os.Stdin)
	newEnvArg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
	newEnvArg = strings.TrimSpace(newEnvArg)
	localEnv[name] = newEnvArg
	Debug("Saving %s value", name)
	godotenv.Write(localEnv, envFile)
}
