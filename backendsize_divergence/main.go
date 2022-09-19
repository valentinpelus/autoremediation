package main

import (
	"aws/autoremediate/backendsize_divergence/cli"
	"context"
	"crypto/tls"
	"flag"
	"net/http"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	//"github.com/valentinpelus/go-package/cli"
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
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	cli.LoadConf(*confPath)

	// Loading kubeconfig file with context
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Init http client
	Client = &http.Client{Transport: tr}

	// Init AMUrl to allow alerts query
	jsonUrl := cli.Conf.QueryURL + "/api/v1/alerts"

	namespace := "ingress-controller-v2"
	podName := "ingress-test2"

	pods, err := clientset.CoreV1().Pods("ingress-controller-v2").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Info().Msgf("List of pod %s", pods)
	log.Info().Msgf("Searching pod %s in namespace %s", podName, namespace)
	// Parsing pod on the namespace to find the right one
	/*for _, podsInfo := range (*pods).Items {
		log.Info().Msgf("Parsing pod %s", podsInfo.Name)
		if podsInfo.Name == podName {
			log.Info().Msgf("Found pod %s in namespace %s", podName, namespace)
		}
	}*/

	for {
		time.Sleep(1 * time.Second)
		log.Info().Msgf("Check ongoing")
		// Querying Alertmanager to check if alert is firing for backend size divergence and proceed to deletion if needed
		//getAlert := getVMAlertBackendSize(jsonUrl, clientset)
		podName, namespace := cli.GetVMAlertBackendSize(jsonUrl)
		if (len(podName) > 0) && (len(namespace) > 0) {
			log.Info().Msgf("Detecting pod %s in namespace %s on divergence", podName, namespace)
			cli.DeletePod(podName, clientset, namespace)
		}
	}
}
