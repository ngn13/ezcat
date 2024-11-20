package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/log"
)

func Error(c *fiber.Ctx, err string) error {
	return c.JSON(fiber.Map{
		"error": err,
	})
}

func ErrorCode(c *fiber.Ctx, code int) error {
	switch code {
	case 400:
		return c.Status(code).JSON(fiber.Map{
			"error": "Bad request",
		})

	case 401:
		return c.Status(code).JSON(fiber.Map{
			"error": "You are not logged in",
		})

	default:
		return c.Status(code).JSON(fiber.Map{
			"error": "Internal error",
		})
	}
}

func ErrorInternal(c *fiber.Ctx, err string) error {
	log.Fail("%s", err)
	return ErrorCode(c, 500)
}
