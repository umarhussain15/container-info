package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/umarhussain15/container-info/internal/handlers/podsearch"
	k8srepo "github.com/umarhussain15/container-info/internal/repositories/k8s-repo"
	"github.com/umarhussain15/container-info/utils"
)

func main() {
	source, err := k8srepo.NewInClusterSource()
	if err != nil {
		log.Fatalln(err)
		return
	}
	handler := podsearch.NewHTTPHandler(source)

	app := fiber.New()
	app.Use(logger.New())
	app.Get("/container-resources", handler.GetContainerInfo)

	host := utils.GetEnvOrDefault("HOST", "0.0.0.0")
	port := utils.GetEnvOrDefault("PORT", "8089")
	err = app.Listen(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalln(err)
		return
	}
}
