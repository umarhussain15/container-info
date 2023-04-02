package ports

import "github.com/umarhussain15/container-info/internal/core/domain"

// K8sRepository port list functionalities a service needs to implement.
type K8sRepository interface {
	SearchPods(podLabel string) ([]domain.ContainerInfo, error)
}
