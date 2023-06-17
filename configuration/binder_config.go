package configuration

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type ServiceApp struct {
	EnvVariable string
	Config      ConfigApp
	Path        string
}

func (s *ServiceApp) Load() {
	if err := envconfig.Process(s.EnvVariable, &s.Config); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
