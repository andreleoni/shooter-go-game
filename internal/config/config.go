package config

import "os"

func IsDebugEnv() bool {
	return os.Getenv("env") == "debug"
}
