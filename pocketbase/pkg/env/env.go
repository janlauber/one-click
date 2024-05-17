package env

import (
	"log"

	"github.com/caarlos0/env/v8"
)

type config struct {
	Local               bool   `env:"LOCAL"`
	LocalKubeConfigFile string `env:"LOCAL_KUBECONFIG_FILE" envDefault:"~/.kube/config"`
	CronTick            string `env:"CRON_TICK" envDefault:"*/1 * * * *"`
}

var Config config

func Init() {
	if err := env.Parse(&Config); err != nil {
		log.Printf("%+v\n", err)
	}

	if Config.Local {
		log.Println("Running in local mode and kubeconfig located at: " + Config.LocalKubeConfigFile)
	}
}
