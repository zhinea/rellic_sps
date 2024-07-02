package containercontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/model/entity"
	"log"
)

type Container struct {
	ID       int    `json:"id" form:"id"`
	Config   string `json:"config" form:"config"`
	IsActive int    `json:"is_active" form:"is_active" gorm:"allowzero"`
}

func Create(c *fiber.Ctx) error {

	createData := new(Container)

	if err := c.BodyParser(createData); err != nil {
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "B1",
			"message":     "Invalid request body",
		})
	}

	// do something with createData
	if err := database.DB.Create(&createData).Error; err != nil {
		log.Println("err: can't insert request log to database", err)

		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "S1",
			"message":     "Server can't save to database.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"ok":          1,
		"status_code": "N",
		"message":     "Success",
	})
}

func Update(c *fiber.Ctx) error {
	updateData := new(Container)

	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":          0,
			"status_code": "B1",
			"message":     "Invalid request body",
		})
	}

	// get user from db
	var container Container
	if err := database.DB.Where("id = ?", updateData.ID).First(&container).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":          0,
			"status_code": "WH0",
			"message":     "Container not found",
		})
	}

	// do something with updateData
	if err := database.DB.Model(&container).Updates(map[string]interface{}{
		"config":    updateData.Config,
		"is_active": updateData.IsActive,
	}).Error; err != nil {
		log.Println("err: can't insert request log to database", err)

		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "SU1",
			"message":     "Server can't save to database.",
		})
	}

	var domainResults []entity.Domain

	database.DB.Model(&entity.Domain{}).
		Where("container_id = ?", container.ID).
		Select("container_id, domain").
		Scan(&domainResults)

	domains := make([]string, len(domainResults))

	for _, result := range domainResults {
		domains = append(domains, "host:"+result.Domain)
	}

	ctx := context.Background()

	errRedis := database.Redis.Del(ctx, domains...)
	if errRedis != nil {
		log.Println(errRedis)
	}

	return c.Status(200).JSON(fiber.Map{
		"ok":          1,
		"status_code": "N",
		"message":     "Success",
	})
}

func Delete(c *fiber.Ctx) error {
	deleteData := new(Container)

	if err := c.BodyParser(deleteData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":          0,
			"status_code": "B1",
			"message":     "Invalid request body",
		})
	}

	// get user from db
	var container Container
	if err := database.DB.Where("id = ?", deleteData.ID).First(&container).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":          0,
			"status_code": "WH0",
			"message":     "Container not found",
		})
	}

	// do something with deleteData
	if err := database.DB.Delete(&container).Error; err != nil {
		log.Println("err: can't insert request log to database", err)

		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "SD1",
			"message":     "Server can't save to database.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"ok":          1,
		"status_code": "N",
		"message":     "Success",
	})
}
