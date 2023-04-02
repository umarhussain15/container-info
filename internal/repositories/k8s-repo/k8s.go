package k8srepo

import (
	"github.com/umarhussain15/container-info/internal/core/domain"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sRepository struct {
	client *kubernetes.Clientset
}

// NewInClusterSource generates an instance of repository which utilizes serviceaccount attached to the pod in cluster.
func NewInClusterSource() (*K8sRepository, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client := kubernetes.NewForConfigOrDie(config)
	return &K8sRepository{
		client: client,
	}, nil
}

// SearchPods searches the cluster for the pods matching the given label search. The label search format is similar to
// kubectl cli format.
func (cluster *K8sRepository) SearchPods(podLabel string) ([]domain.ContainerInfo, error) {
	options := metav1.ListOptions{
		LabelSelector: podLabel,
	}
	list, err := cluster.client.CoreV1().Pods("").List(context.Background(), options)
	if err != nil {
		return nil, err
	}
	var containers []domain.ContainerInfo
	for _, pod := range list.Items {
		for _, container := range pod.Spec.Containers {
			info := domain.ContainerInfo{
				ContainerName: container.Name,
				PodName:       pod.Name,
				Namespace:     pod.Namespace,
				MemoryRequest: container.Resources.Requests.Memory().String(),
				MemoryLimit:   container.Resources.Limits.Memory().String(),
				CPURequest:    container.Resources.Requests.Cpu().String(),
				CPULimit:      container.Resources.Limits.Cpu().String(),
			}
			containers = append(containers, info)
		}
	}
	return containers, nil
}
