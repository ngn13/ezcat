package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/config"
)

func GET_info(c *fiber.Ctx) error {
	conf := c.Locals("config").(*config.Struct)

	return c.JSON(&fiber.Map{
		"version":  conf.Version,
		"megamind": conf.Megamind,
	})
}
