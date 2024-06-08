package routes 

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/util"
)

func GET_logout(c *fiber.Ctx) error {
  TokenDel(util.GetToken(c))
  return c.JSON(fiber.Map{})
}
