package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

type Response struct {
	Status string `json:"status"`
	Data   []Data `json:"data"`
}
type Labels struct {
	AdminAlert              string `json:"admin_alert"`
	Alertgroup              string `json:"alertgroup"`
	Alertname               string `json:"alertname"`
	ClusterName             string `json:"cluster_name"`
	Horizontalpodautoscaler string `json:"horizontalpodautoscaler"`
	MetricName              string `json:"metric_name"`
	Namespace               string `json:"namespace"`
	Pod                     string `json:"pod"`
	NmcbackAlert            string `json:"nmcback_alert"`
	OnCall                  string `json:"on_call"`
	Severity                string `json:"severity"`
	Team                    string `json:"team"`
}

type Data struct {
	Labels Labels `json:"labels,omitempty"`
}

func getVMAlertBackendSize(server string) bool {
	res, err := http.Get(server)
	if err != nil {
		log.Fatal().Msgf("Error in GET request %s ", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Msgf("Error in reading body %s ", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Can't unmarshal JSON")
	}

	alertFiring := false

	for _, alerts := range response.Data {
		if (alerts.Labels.Alertname == "haproxyBackendSizeDivergence") && (len(alerts.Labels.Pod) > 0) {
			log.Info().Msgf("Alert is firing %s on pod %s ", alerts.Labels.Alertname, alerts.Labels.Pod)
			//podName := alerts.Labels.Pod
			alertFiring = true
		} else {
			log.Info().Msgf("No pod in state of backendsize divergence")
			alertFiring = false
		}
	}
	return alertFiring
}

func main() {

	// Init of zerolog library
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Init http client
	Client = &http.Client{
		Timeout: time.Second * 2,
	}

	// Querying Alertmanager to check if alert is firing for backend size divergence
	server := "https://alertmanager.staging.6cloud.fr/api/v1/alerts"

	//log.Info().Msgf("Alert is firing %s ", response)
	getVMAlertBackendSize(server)
}
