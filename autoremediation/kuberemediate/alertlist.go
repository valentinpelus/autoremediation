package kuberemediate

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type AlertList struct {
	Alertlist []struct {
		Alertname string `yaml:"alertname"`
		Enabled   bool   `yaml:"enabled"`
		Action    string `yaml:action`
	} `yaml:"alertlist"`
}

var Alerts AlertList

func LoadAlertList(confPath string) {
	yamlFile, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't read kuberemediate config file")
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &Alerts)
	if err != nil {
		log.Fatal().Err(err).Msg("Unmarshal error")
		os.Exit(1)
	}
}
