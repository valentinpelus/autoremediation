package kuberemediate

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/rs/zerolog/log"
)

func DeletePod(podName string, clientset *kubernetes.Clientset, namespace string) bool {

	if CheckPodPresent(podName, clientset, namespace) {
		log.Info().Msgf("Deleting pod %s", podName)
		if err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{}); err != nil {
			log.Info().Msgf("Error in deletion of pod %s ", podName)
			panic(err)
		}
		return true
	}
	return false
}

func CheckPodPresent(podName string, clientset *kubernetes.Clientset, namespace string) bool {

	// Listing Pods from chosen namespace
	pods, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: ""})
	log.Info().Msgf("Searching pod %s in namespace %s", podName, namespace)
	// Checking if we have more than one pod on our namespace, if it's the case then we can proceed
	if checkQuotaPodOnNamespace(clientset, namespace) {
		log.Info().Msgf("More than 2 pod on namespace %s can proceed to actions", namespace)
		// Parsing pod on the namespace to find the right one
		for _, podsInfo := range (*pods).Items {
			if (podsInfo.Name == podName) && (podsInfo.Status.Phase == "Running") {
				log.Info().Msgf("Found pod %s in namespace %s in status %s", podName, namespace, podsInfo.Status.Phase)
				return true
			}
		}
	} else {
		log.Error().Msgf("Error during deletion, we have less than 2 pods on cluster, therefore no action will be taken ")
	}
	return false
}

func checkQuotaPodOnNamespace(clientset *kubernetes.Clientset, namespace string) bool {

	// Using this func to check if we have more than one pod on our namespace before taking any action, avoiding creating chain reaction
	pods, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: ""})
	log.Info().Msgf("Checking number of pod on namespace %d before taking actions", namespace)
	return len(pods.Items) >= 2
}
