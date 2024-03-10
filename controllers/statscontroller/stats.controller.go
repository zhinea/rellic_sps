package statscontroller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/database"
)

func ContainerRequestLogs(c *fiber.Ctx) error {
	var result []map[string]interface{}

	err := database.DB.
		Raw("SELECT created_at, COUNT(*) AS usage_count, container_id FROM request_logs WHERE created_at >= NOW() - INTERVAL 5 MINUTE GROUP BY container_id").
		Scan(&result).
		Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":          0,
			"status_code": "ECS21",
			"message":     "Failed to fetch container request logs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":          1,
		"status_code": "N",
		"result":      result,
	})
}
