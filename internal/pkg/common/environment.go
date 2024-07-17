package common

import (
	"os"
	"fmt"

	"github.com/joho/godotenv"
)

var localEnv map[string]string

func GetEnvArg(name string, defaultValue ...string) (string) {
	value string
	if val, ok := localEnv[name]; ok {
		value = val
	} else if envVal, exists := os.LookupEnv(name); exists {
		value = envVal
	} else if len(defaultValue) > 0 {
		value = defaultValue[0]
	} else {
		Error(fmt.Sprintf("Env var '%s' is not present", name))
		os.Exit(1)
	}
	return value
}

func LoadEnvironment()( map[string]string) {
	Debug("Checking environment packages and tools")
	CommandExists("git")
	Debug("Loading environment variables from ~/.gmenv file")

	var localEnv map[string]string
	localEnv, err := godotenv.Read()

	err := godotenv.Load(string(os.Getenv("HOME")) + "/.gmenv")
	if err != nil {
		Error("Error loading ~/.gmenv file")
		os.Exit(1)
	}

}
