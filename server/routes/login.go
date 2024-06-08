package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/global"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

func PUT_login(c *fiber.Ctx) error {
  var data map[string]string
  if err := c.BodyParser(&data); err != nil {
    return util.ErrorCode(c, 400)
  }

  if data["password"] == "" {
    return util.ErrorCode(c, 400)
  } 

  if global.CONFIG_PASSWORD != data["password"] {
    return util.Error(c, "Invalid password")
  }

  log.Info("New login from %s", c.IP())
  return c.JSON(fiber.Map{
    "token": TokenNew(),
  })
}
