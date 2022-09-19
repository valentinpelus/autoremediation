package cli

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/rs/zerolog/log"
	//"github.com/valentinpelus/go-package/cli"
)

func DeletePod(podName string, clientset *kubernetes.Clientset, namespace string) bool {

	podExists := CheckPodPresent(podName, clientset, namespace)

	/*if pod, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{}); err != nil {
		log.Fatal().Msgf("Error in select %s ", err)
		println(pod)
		panic(err.Error())
	}*/

	if podExists {
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
	// Namespace on this alert is set at static on ingress-controller-v2
	//namespace := "ingress-controller-v2"

	// Listing Pods from chosen namespace
	pods, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: ""})
	log.Info().Msgf("List of pod %d", pods)
	log.Info().Msgf("Searching pod %s in namespace %s", podName, namespace)
	// Parsing pod on the namespace to find the right one
	for _, podsInfo := range (*pods).Items {
		log.Info().Msgf("Parsing pod %s", podsInfo.Name)
		if podsInfo.Name == podName {
			log.Info().Msgf("Found pod %s in namespace %s", podName, namespace)
			return true
		}
	}
	return false
}
