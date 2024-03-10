package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/controllers/containercontroller"
	"github.com/zhinea/sps/controllers/domaincontroller"
	"github.com/zhinea/sps/controllers/statscontroller"
)

func SysRouteInit(router fiber.Router) {

	v1 := router.Group("/v1")

	container := v1.Group("/containers")

	container.Post("/", containercontroller.Create)
	container.Put("/", containercontroller.Update)
	container.Delete("/", containercontroller.Delete)

	domain := v1.Group("/domains")

	domain.Post("/", domaincontroller.Create)
	domain.Put("/", domaincontroller.Update)
	domain.Delete("/", domaincontroller.Delete)

	stats := v1.Group("/stats")

	stats.Get("/get_containers_sden", statscontroller.ContainerRequestLogs)
}
