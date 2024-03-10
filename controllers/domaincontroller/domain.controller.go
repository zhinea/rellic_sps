package domaincontroller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/database"
	"log"
)

type Domain struct {
	ID          string `json:"id" form:"id"`
	ContainerID int    `json:"container_id" form:"container_id"`
	Domain      string `json:"domain" form:"domain"`
}

func Create(c *fiber.Ctx) error {
	createData := new(Domain)

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

func deleteCache(ctx context.Context, domain string) error {
	// delete cache
	return database.Redis.Del(ctx, "host:"+domain).Err()
}

func Update(c *fiber.Ctx) error {
	updateData := new(Domain)

	if err := c.BodyParser(updateData); err != nil {
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "B1",
			"message":     "Invalid request body",
		})
	}

	// get user from db
	var domain Domain
	if err := database.DB.Where("id = ?", updateData.ID).First(&domain).Error; err != nil {
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "WH0",
			"message":     "Domain not found",
		})
	}

	// do something with updateData
	if err := database.DB.Model(&domain).Updates(map[string]interface{}{
		"domain": updateData.Domain,
	}).Error; err != nil {
		log.Println("err: can't insert request log to database", err)
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "SU1",
			"message":     "Server can't save to database.",
		})
	}

	// delete cache
	if err := deleteCache(c.Context(), domain.Domain); err != nil {
		log.Println("err: can't delete cache", err)
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "SC2",
			"message":     "Server can't delete cache.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"ok":          1,
		"status_code": "N1",
		"message":     "Success",
	})
}

func Delete(c *fiber.Ctx) error {
	deleteData := new(Domain)

	if err := c.BodyParser(deleteData); err != nil {
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "B1",
			"message":     "Invalid request body",
		})
	}

	// get user from db
	var domain Domain
	if err := database.DB.Where("id = ?", deleteData.ID).First(&domain).Error; err != nil {
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "WH0",
			"message":     "Domain not found",
		})
	}

	// do something with deleteData
	if err := database.DB.Delete(&domain).Error; err != nil {
		log.Println("err: can't insert request log to database", err)
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "SD1",
			"message":     "Server can't save to database.",
		})
	}

	// delete cache
	if err := deleteCache(c.Context(), domain.Domain); err != nil {
		log.Println("err: can't delete cache", err)
		return c.Status(200).JSON(fiber.Map{
			"ok":          0,
			"status_code": "SC2",
			"message":     "Server can't delete cache.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"ok":          1,
		"status_code": "N",
		"message":     "Success",
	})
}
