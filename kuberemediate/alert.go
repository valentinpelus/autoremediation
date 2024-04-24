package kuberemediate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Response struct {
	Status string `json:"status"`
	Data   []Data `json:"data"`
}
type Labels struct {
	AdminAlert  string `json:"admin_alert"`
	Alertgroup  string `json:"alertgroup"`
	AlertName   string `json:"alertname"`
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Pod         string `json:"pod"`
}

type Data struct {
	Labels Labels `json:"labels,omitempty"`
}

var (
	Client HTTPClient
)

var alertPodExtractList []string

func GetVMAlertMatch(server string, ListSupportedAlert []string) (string, string, string) {
	// Initialisation of GET request
	res, err := http.Get(server)
	if err != nil {
		log.Fatal().Msgf("Error in GET request %s ", err)
	}
	// Closing request
	defer res.Body.Close()

	// Reading Body content
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Msgf("Error in reading body %s ", err)
	}

	// Serialising return of Body into JSON
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal().Msgf("Error in reading body %s ", err)
	}

	for _, list := range ListSupportedAlert {
		fmt.Printf("%+v\n", list)
	}

	// Parsing Json return to match Alertname with the slice sent to the function
	for _, alerts := range response.Data {
		for _, supportedAlert := range ListSupportedAlert {
			if alerts.Labels.AlertName == supportedAlert && len(alerts.Labels.Pod) > 0 {
				alertPodExtractList = append(alertPodExtractList, alerts.Labels.Pod)
				log.Info().Msgf("Alert %s is firing on pod %s deletion ongoing", alerts.Labels.AlertName, alerts.Labels.Pod)
				podName := alerts.Labels.Pod
				namespace := alerts.Labels.Namespace
				alertName := alerts.Labels.AlertName
				log.Info().Msgf("alert.go Podname : %s Namespace : %s", podName, namespace)
				//Proceeding to the deletion of pod if alert is firing
				return podName, namespace, alertName
			} else {
				log.Info().Msgf("No pod in state of backendsize divergence")
				continue
			}
		}
	}
	return "", "", ""
}
