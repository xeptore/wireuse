package env

import (
	"os"

	"github.com/rs/zerolog/log"
)

func MustGet(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || len(v) == 0 {
		log.Fatal().Str("key", key).Msg("environment variable must be set with a value")
	}

	return v
}
