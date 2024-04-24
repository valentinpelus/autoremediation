package kuberemediate

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Alert struct {
	AlertName string `yaml:"alertname"`
	Enabled   bool   `yaml:"enabled"`
}

type AlertList struct {
	AlertList []Alert `yaml:"alertlist"`
}

var Alerts AlertList

var EnabledAlertList []string

func LoadConfKubeAlert(confPath string) {
	yamlFile, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't read kuberemediate config file")
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &Alerts)
	if err != nil {
		log.Fatal().Err(err).Msg("Unmarshal error on alert list in configuration")
		os.Exit(1)
	}

	confAlertList := Alerts.AlertList
	for _, allAlerts := range confAlertList {
		if allAlerts.Enabled {
			fmt.Printf("%+v\n", allAlerts.AlertName)
			EnabledAlertList = append(EnabledAlertList, allAlerts.AlertName)
		}
	}
}
