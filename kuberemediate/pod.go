package kuberemediate

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/rs/zerolog/log"
)

var podLabelTarget string

func DeletePod(podName string, clientset *kubernetes.Clientset, namespace string, podAmmount int) bool {

	if CheckPodPresent(podName, clientset, namespace, podAmmount) {
		log.Info().Msgf("Deleting pod %s", podName)
		if err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{}); err != nil {
			log.Info().Msgf("Error in deletion of pod %s ", podName)
			panic(err)
		}
		return true
	}
	return false
}

func GetPod(podName string, clientset *kubernetes.Clientset, namespace string) (map[string]string, error) {
	// Get pod by it's name and check if it's present in the namespace, it will help to target the required project's pods with it's label
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Error().Msgf("Error in getting pod %s from namespace %s", podName, namespace)
		panic(err)
	}
	podMap := make(map[string]string)
	podMap["Name"] = pod.Name
	podMap["Namespace"] = pod.Namespace
	podMap["LabelInstance"] = pod.Labels["app.kubernetes.io/instance"]
	podMap["LabelName"] = pod.Labels["app.kubernetes.io/name"]
	podMap["LabelProject"] = pod.Labels["project"]
	return podMap, nil
}

func CheckPodPresent(podName string, clientset *kubernetes.Clientset, namespace string, podAmmount int) bool {

	// Get pod by it's name and check if it's present in the namespace, it will help to target the required project's pods with it's label
	podQuery, err := GetPod(podName, clientset, namespace)
	if err != nil {
		log.Error().Msgf("Error in getting pod %s from namespace %s", podName, namespace)
		panic(err)
	}

	// We will check the labels of our pod to find the right target to list all pods concerning the same project
	if _, ok := podQuery["LabelProject"]; ok {
		// If the label project exist we retrieve it and set it as target"
		podLabelTarget = podQuery["LabelProject"]
	} else {
		// If there is no label project we set the label app.kubernetes.io/name as target
		podLabelTarget = podQuery["app.kubernetes.io/name"]
	}

	// Listing Pods from chosen namespace, targeting the right label
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: podLabelTarget})
	if err != nil {
		log.Error().Msgf("Error in listing pods from namespace %s", namespace)
		panic(err)
	}
	log.Info().Msgf("Searching pod %s in namespace %s", podName, namespace)
	// Using func checkQuotaPod to check the ammount of healthy pod on our project
	if checkQuotaPod(clientset, namespace, podLabelTarget, podAmmount) {
		log.Info().Msgf("More than 2 pod on namespace %s can proceed to actions", namespace)
		// Making sure we are targeting running pod and not backoff/restarting one
		for _, podsInfo := range (pods).Items {
			if (podsInfo.Name == podName) && (podsInfo.Status.Phase == "Running") {
				log.Info().Msgf("Found pod %s in namespace %s in status %s", podName, namespace, podsInfo.Status.Phase)
				return true
			}
		}
	} else {
		log.Warn().Msgf("Warning issued during deletion, the ratio between unhealthy pod and healthy one is not optimal, therefore no action will be taken. Backing off")
	}
	return false
}

func checkQuotaPod(clientset *kubernetes.Clientset, namespace string, podLabelTarget string, podAmmount int) bool {

	// Using this func to check if we have more than one pod on our namespace before taking any action, avoiding creating chain reaction
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: podLabelTarget})
	if err != nil {
		log.Error().Msgf("Error in getting number pods from namespace %s", namespace)
		panic(err)
	}
	log.Info().Msgf("Checking number of pod on namespace %s before taking actions", namespace)

	// If we have more than the ammount of pod returned by the alert - 2 and if we have at minimum 2 pod on the project running, then we can proceed
	// It aims to avoid deleting all pods of the same project directly
	return len(pods.Items) >= podAmmount-2 && len(pods.Items) >= 2
}
