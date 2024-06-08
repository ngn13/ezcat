package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/global"
)

const VERSION = "2.3"

func GET_info(c *fiber.Ctx) error {
  return c.JSON(&fiber.Map{
    "version": VERSION,
    "megamind": global.CONFIG_MEGAMIND,
  })
}
