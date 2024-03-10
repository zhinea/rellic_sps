package statscontroller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/database"
	"time"
)

func ContainerRequestLogs(c *fiber.Ctx) error {
	var result []map[string]interface{}

	pastTime := time.Now().Add(-5 * time.Minute)

	err := database.DB.
		Raw("SELECT created_at, COUNT(*) AS usage_count, container_id FROM request_logs WHERE created_at >= ? GROUP BY container_id", pastTime).
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
