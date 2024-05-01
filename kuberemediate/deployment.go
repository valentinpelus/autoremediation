package kuberemediate

import (
	//"context"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//"github.com/rs/zerolog/log"
)

func EnrichAlerts(clientset *kubernetes.Clientset, namespace string) {

	for _, alert := range EnabledAlertList {
		if alert[1] == "delete" {
			//DescribeDeployment(alert[0], clientset, namespace)
		}
	}
}

/* func DescribeDeployment(deploymentName string, clientset *kubernetes.Clientset, namespace string) bool {

	if CheckPodPresent(deploymentName, clientset, namespace, podAmmount) {
		log.Info().Msgf("Deleting pod %s", deploymentName)
		if err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{}); err != nil {
			log.Info().Msgf("Error in deletion of pod %s ", deploymentName)
			panic(err)
		}
		return true
	}
	return false
}
*/
