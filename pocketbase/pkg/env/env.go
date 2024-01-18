package env

import (
	"log"

	"github.com/caarlos0/env/v8"
)

type config struct {
	Local             bool   `env:"LOCAL"`
	DefaultRolloutDir string `env:"DEFAULT_ROLLOUT_DIR" envDefault:".rollouts"`
}

var Config config

func Init() {
	if err := env.Parse(&Config); err != nil {
		log.Printf("%+v\n", err)
	}

	if Config.Local {
		log.Println("Running in local mode")
	}
}
