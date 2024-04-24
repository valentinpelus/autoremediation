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

<<<<<<< HEAD:autoremediation/main.go
	kuberemediate.LoadConf(*confPath)
=======
	// Loading configuration files
	kuberemediate.LoadConfKube(*confPath)
	kuberemediate.LoadConfKubeAlert(*confPath)
>>>>>>> e51e388 (feat(alertlist): add feature to allow alertlist and enable, disabled them through configuration):main.go

	// Loading kubeconfig file with context
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
<<<<<<< HEAD:autoremediation/main.go
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
=======
	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(kube_config)
>>>>>>> e51e388 (feat(alertlist): add feature to allow alertlist and enable, disabled them through configuration):main.go
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
<<<<<<< HEAD:autoremediation/main.go
		time.Sleep(60 * time.Second)
		log.Info().Msgf("Check ongoing")
		// Querying Alertmanager to check if alert is firing for backend size divergence and proceed to deletion if needed
		podName, namespace := kuberemediate.GetVMAlertBackendSize(jsonUrl)
=======
		time.Sleep(10 * time.Second)
		log.Info().Msgf("Check ongoing")
		// Querying Alertmanager to check if alert is firing for backend size divergence and proceed to deletion if needed
		podName, namespace, alertname := kuberemediate.GetVMAlertMatch(jsonUrl, ListSupportedAlert)
>>>>>>> e51e388 (feat(alertlist): add feature to allow alertlist and enable, disabled them through configuration):main.go
		if (len(podName) > 0) && (len(namespace) > 0) {
			log.Info().Msgf("Detecting pod %s in namespace %s on divergence", podName, namespace)
			kuberemediate.DeletePod(podName, clientset, namespace)
		}
	}
}
