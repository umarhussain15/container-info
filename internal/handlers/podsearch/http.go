package podsearch

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/umarhussain15/container-info/internal/core/ports"
)

type HTTPHandler struct {
	k8s ports.K8sRepository
}

// NewHTTPHandler provides HTTPHandler instance which has handlers for gofiber.
func NewHTTPHandler(k8s ports.K8sRepository) *HTTPHandler {
	return &HTTPHandler{
		k8s: k8s,
	}
}

// GetContainerInfo extracts the query for the label search and call the search on the k8s repository for matching pods.
func (receiver *HTTPHandler) GetContainerInfo(c *fiber.Ctx) error {
	query := c.Query("pod-label")

	pods, err := receiver.k8s.SearchPods(query)
	if err != nil {
		return err
	}
	if pods == nil || len(pods) < 1 {
		return c.Status(404).JSON(fiber.Map{
			"message": fmt.Sprintf("No pod found for %s", query),
		})
	}
	return c.JSON(&pods)
}
