package kuberemediate

import (
	"encoding/json"
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

func GetVMAlertMatch(server string, ListSupportedAlert [][]string) (string, string, string, string) {
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

	// Parsing Json return to match Alertname with the slice sent to the function
	for _, alerts := range response.Data {
		alertName := alerts.Labels.AlertName
		// Parsing all our supported alerts, if we find a match we append it to our slice then return it at the end
		for i := range ListSupportedAlert {
			if alertName == ListSupportedAlert[i][0] && len(alerts.Labels.Pod) > 0 {
				podName := alerts.Labels.Pod
				namespace := alerts.Labels.Namespace
				alertAction := ListSupportedAlert[i][1]
				alertPodExtractList = append(alertPodExtractList, podName)
				log.Info().Msgf("Alert %s is firing", alertName)
				if podName != "" && namespace != "" {
					log.Info().Msgf("alert.go Podname : %s Namespace : %s", podName, namespace)
				}
				//Proceeding to the deletion of pod if alert is firing
				return podName, namespace, alertName, alertAction
			} else {
				log.Info().Msgf("No alerts matching any enabled rules in remediate, continuing...")
				continue
			}
		}
	}
	return "", "", "", ""
}
