package kuberemediate

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Alert struct {
	AlertName string `yaml:"alertname"`
	Enabled   bool   `yaml:"enabled"`
	Action    string `yaml:"action"`
}

type AlertList struct {
	AlertList []Alert `yaml:"alertRulesList"`
}

var Alerts AlertList
var EnabledAlertList [][]string

func LoadConfAlert(confPath string) {
	yamlFile, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't read kuberemediate alerts config file")
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &Alerts)
	if err != nil {
		log.Fatal().Err(err).Msg("Unmarshal error on alert list in configuration")
		os.Exit(1)
	}

	confAlertList := Alerts.AlertList
	for i := range confAlertList {
		if confAlertList[i].Enabled {
			EnabledAlertList = append(EnabledAlertList, []string{confAlertList[i].AlertName, confAlertList[i].Action})
		}
	}
}
