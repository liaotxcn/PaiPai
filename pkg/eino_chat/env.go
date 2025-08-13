package eino_chat

import (
	"log"
	"os"
)

func MustHasEnvs(envs ...string) {
	for _, env := range envs {
		if os.Getenv(env) == "" {
			log.Fatalf("❌ [ERROR] env [%s] is required, but is not set now, please check your .env file", env)
		}
	}
}
