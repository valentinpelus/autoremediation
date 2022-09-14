package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	//appsv1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/valentinpelus/go-package/cli"
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

func getPodDivergence(podName string, clientset *kubernetes.Clientset) bool {
	// Namespace on this alert is set at static on ingress-controller-v2
	namespace := "ingress-controller-v2"

	// Listing Pods from chosen namespace
	pods, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: ""})
	log.Info().Msgf("Searching pod %s in namespace %s", podName, namespace)
	// Parsing pod on the namespace to find the one in backend size divergence and to be sure it still exists on the namespace
	for _, podsInfo := range (*pods).Items {
		log.Info().Msgf("Parsing pod %s", podsInfo.Name)
		if podsInfo.Name == podName {
			log.Info().Msgf("Found pod %s in namespace %s", podName, namespace)
			return true
			break
		}
	}
	return false
}

func deletePodDivergence(podName string, clientset *kubernetes.Clientset) bool {

	podExists := getPodDivergence(podName, clientset)

	if podExists {
		log.Info().Msgf("Deleting pod %s in state of backend size divergence", podName)
		if err := clientset.CoreV1().Pods("webserver").Delete(context.TODO(), podName, metav1.DeleteOptions{}); err != nil {
			log.Info().Msgf("Error in deletion of pod %s ", podName)
			panic(err)
			return false
		}
	}
	return true
}

func getVMAlertBackendSize(server string, clientset *kubernetes.Clientset) bool {
	// Initialisation of GET request
	res, err := http.Get(server)
	if err != nil {
		log.Fatal().Msgf("Error in GET request %s ", err)
	}

	// Closing request
	defer res.Body.Close()

	// Reading Body content
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Msgf("Error in reading body %s ", err)
	}

	// Serialising return of Body into JSON
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal().Msgf("Error in reading body %s ", err)
	}

	// Init var alertFiring with false by default
	alertFiring := false

	// Parsing Json return to match Alertname with haproxyBackendSizeDivergence
	for _, alerts := range response.Data {
		if (alerts.Labels.Alertname == "haproxyBackendSizeDivergence") && (len(alerts.Labels.Pod) > 0) {
			log.Info().Msgf("Alert %s is firing on pod %s deletion ongoing", alerts.Labels.Alertname, alerts.Labels.Pod)
			podName := alerts.Labels.Pod
			//Proceeding to the deletion of pod if alert is firing
			deletePodDivergence(podName, clientset)
			alertFiring = true
			break
		} else {
			log.Info().Msgf("No pod in state of backendsize divergence")
			alertFiring = false
		}
	}
	return alertFiring
}

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

	for {
		time.Sleep(1 * time.Second)

		// Querying Alertmanager to check if alert is firing for backend size divergence and proceed to deletion if needed
		getAlert := getVMAlertBackendSize(jsonUrl, clientset)
		log.Info().Msgf("Check ongoing on %s ", getAlert)
	}
}
