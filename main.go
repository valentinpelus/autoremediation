package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"remediate/kuberemediate"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func main() {

	// Init of zerolog library and config
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	confPath := flag.String("conf", "config.yaml", "Config path")
	debug := flag.Bool("debug", false, "sets log level to debug")
	//tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Loading configuration files
	kuberemediate.LoadConfKube(*confPath)
	kuberemediate.LoadConfAlert(*confPath)

	// Loading kubeconfig file with context
	kube_config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(kube_config)
	if err != nil {
		panic(err.Error())
	}

	// Init http client
	Client = &http.Client{}

	// Init AMUrl to allow alerts query
	jsonUrl := kuberemediate.Conf.QueryURL + "/api/v1/alerts"

	ListSupportedAlert := kuberemediate.EnabledAlertList

	fmt.Println(ListSupportedAlert)

	for {
		time.Sleep(20 * time.Second)
		log.Info().Msgf("Check ongoing")
		// Querying Alertmanager to check if alert is firing for backend size divergence and proceed to deletion if needed
		alertPodExtractList := kuberemediate.GetVMAlertMatch(jsonUrl, ListSupportedAlert)
		fmt.Println("Print of return alertPodExtractList : ", alertPodExtractList)
		for i := range alertPodExtractList {
			podName := alertPodExtractList[i][0]
			namespace := alertPodExtractList[i][1]
			alertAction := alertPodExtractList[i][2]
			alertName := alertPodExtractList[i][3]
			if (len(podName) > 0) && (len(namespace) > 0) {
				log.Info().Msgf("Detecting pod %s in namespace %s on divergence", podName, namespace)

				// Parse returned alertPodExtractList to determine which action should be done with remediate
				switch alertAction {
				case "deletePod":
					log.Info().Msgf("Delete pod %s in namespace %s on divergence", podName, namespace)
					kuberemediate.DeletePod(podName, clientset, namespace, len(alertPodExtractList))
					time.Sleep(5 * time.Second)
					postMessageSlack(alertName, namespace, confPath)
				case "enrichAlert":
					//kuberemediate.DescribeDeployment(podName, clientset, namespace)

				}
			}
		}
	}
}
