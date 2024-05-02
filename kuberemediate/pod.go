package kuberemediate

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/rs/zerolog/log"
)

var podLabelTarget string

func DeletePod(podInfo map[string]interface{}, clientset *kubernetes.Clientset) bool {

	gracePeriod := int64(0)
	if CheckPodPresent(podInfo, clientset) {
		log.Info().Msgf("Deleting pod %s", podInfo["podName"])
		if err := clientset.CoreV1().Pods(podInfo["namespace"].(string)).Delete(context.TODO(), podInfo["podName"].(string), metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod}); err != nil {
			log.Info().Msgf("Error in deletion of pod %s", podInfo["podName"])
			panic(err)
		}
		return true
	}
	return false
}

func GetPod(podName string, namespace string, clientset *kubernetes.Clientset) (map[string]string, error) {

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
	fmt.Println(podMap)
	return podMap, nil
}

func CheckPodPresent(podInfo map[string]interface{}, clientset *kubernetes.Clientset) bool {

	// Get pod by it's name and check if it's present in the namespace, it will help to target the required project's pods with it's label
	podQuery, err := GetPod(podInfo["podName"].(string), podInfo["namespace"].(string), clientset)
	if err != nil {
		log.Error().Msgf("Error in getting pod %s from namespace %s", podInfo["podName"], podInfo["namespace"])
		panic(err)
	}

	// We will check the labels of our pod to find the right target to list all pods concerning the same project
	if _, ok := podQuery["LabelProject"]; ok {
		// If the label project exist we retrieve it and set it as target"
		podLabelTarget = "project=" + podQuery["LabelProject"]
	} else {
		// If there is no label project we set the label app.kubernetes.io/name as target
		podLabelTarget = "app.kubernetes.io/name=" + podQuery["app.kubernetes.io/name"]
	}

	fmt.Println("Pod Label Target : ", podLabelTarget)

	// Listing Pods from chosen namespace, targeting the right label
	pods, err := clientset.CoreV1().Pods(podInfo["namespace"].(string)).List(context.TODO(), metav1.ListOptions{LabelSelector: podLabelTarget})
	if err != nil {
		log.Error().Msgf("Error in listing pods from namespace %s", podInfo["namespace"])
		panic(err)
	}
	log.Info().Msgf("Searching pod %s in namespace %s", podInfo["podName"], podInfo["namespace"])
	// Using func checkQuotaPod to check the ammount of healthy pod on our project
	podCount, _ := podInfo["podCount"].(int)
	if checkQuotaPod(podInfo["namespace"].(string), podLabelTarget, podCount, clientset) {
		log.Info().Msgf("More than 2 pod on namespace %s can proceed to actions", podInfo["namespace"])
		// Making sure we are targeting running pod and not backoff/restarting one
		for _, podsList := range (pods).Items {
			if (podsList.Name == podInfo["podName"]) && (podsList.Status.Phase == "Running") {
				log.Info().Msgf("Found pod %s in namespace %s in status %s", podInfo["podName"], podInfo["namespace"], podsList.Status.Phase)
				return true
			}
		}
	} else {
		log.Warn().Msgf("Warning issued during deletion, the ratio between unhealthy pod and healthy one is not optimal, therefore no action will be taken. Backing off")
	}
	return false
}

func checkQuotaPod(namespace string, podLabelTarget string, podAmmount int, clientset *kubernetes.Clientset) bool {

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
