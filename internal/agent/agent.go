package agent

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ScaleConfig holds scaling rules
type ScaleConfig struct {
	Namespace    string
	Deployment   string
	MaxReplicas  int
	MinReplicas  int
	CPUThreshold int
}

// ScaleApplication scales up/down based on CPU usage
func ScaleApplication(clientset *kubernetes.Clientset, config ScaleConfig) {
	for {
		deployment, err := clientset.AppsV1().Deployments(config.Namespace).Get(context.TODO(), config.Deployment, metav1.GetOptions{})
		if err != nil {
			log.Printf("Error fetching deployment: %v\n", err)
			continue
		}

		replicas := *deployment.Spec.Replicas

		// Simulated logic for scaling
		if replicas < int32(config.MaxReplicas) {
			replicas++
			log.Println("Scaling up...")
		} else {
			replicas = 0 // Remove all resources
			log.Println("Scaling down: Removing all replicas")
		}

		// Apply scaling decision
		deployment.Spec.Replicas = &replicas
		_, err = clientset.AppsV1().Deployments(config.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("Error updating deployment: %v\n", err)
		}

		time.Sleep(30 * time.Second) // Adjust polling interval
	}
}

// StartAgent initializes the scaling agent
func StartAgent() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error loading Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	scaleConfig := ScaleConfig{
		Namespace:    "default",
		Deployment:   "my-app",
		MaxReplicas:  5,
		MinReplicas:  1,
		CPUThreshold: 70,
	}

	log.Println("Starting Kubernetes Scaling Agent...")
	ScaleApplication(clientset, scaleConfig)
}
