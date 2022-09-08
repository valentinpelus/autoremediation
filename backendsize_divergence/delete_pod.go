package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"path/filepath"

	//appsv1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

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

var (
	Client HTTPClient
)

func getPodDivergence(podName string, clientset *kubernetes.Clientset) bool {
	// Namespace on this alert is set at static on ingress-controller-v2
	namespace := "webserver"

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
			podName = "web-test-56c56dc94b-csktl"

			//Proceeding to the deletion of pod if alert is firing
			deletePodDivergence(podName, clientset)
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

	// Loading kubeconfig file with context
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Init http client
	Client = &http.Client{}

	// Querying Alertmanager to check if alert is firing for backend size divergence and proceed to deletion if needed
	server := "https://alertmanager.staging.6cloud.fr/api/v1/alerts"
	getAlert := getVMAlertBackendSize(server, clientset)
	log.Info().Msgf("Alert is firing %s ", getAlert)
}

func int32Ptr(i int32) *int32 { return &i }
