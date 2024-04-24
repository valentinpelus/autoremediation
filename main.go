package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/valentinpelus/go-package/kuberemediate"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/rs/zerolog"
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
	Alertname   string `json:"alertname"`
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

	kuberemediate.LoadConfKube(*confPath)

	// Loading kubeconfig file with context
	kube_config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(kube_config)
	if err != nil {
		panic(err.Error())
	}

	// Init http client
	Client = &http.Client{}

	// Init AMUrl to allow alerts query
	jsonUrl := kuberemediate.Conf.QueryURL + "/api/v1/alerts"

	for {
		time.Sleep(120 * time.Second)
		log.Info().Msgf("Check ongoing")
		// Querying Alertmanager to check if alert is firing for backend size divergence and proceed to deletion if needed
		podName, namespace, alertname := kuberemediate.GetVMAlertBackendSize(jsonUrl)
		if (len(podName) > 0) && (len(namespace) > 0) {
			log.Info().Msgf("Detecting pod %s in namespace %s on divergence", podName, namespace)
			kuberemediate.DeletePod(podName, clientset, namespace)
			postMessageSlack(alertname, namespace, confPath)
		}
	}
}
