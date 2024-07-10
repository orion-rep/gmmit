package common

import (
	"os"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func GetEnvArg(name string, defaultValue ...string) (string) {
	value, exists := os.LookupEnv(name)
    if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		} else {
			Warning(fmt.Sprintf("Env var '%s' is not present", name))
			os.Exit(1)
		}
	}
	return value
}

func LoadEnvironment() {
	Debug("Loading environment variables from ~/.gmenv file")
	err := godotenv.Load(string(os.Getenv("HOME")) + "/.gmenv")
	if err != nil {
		log.Fatal("Error loading ~/.gmenv file")
	}

}
