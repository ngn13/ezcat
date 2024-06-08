package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/payload"
)

func GET_stage(c *fiber.Ctx) error {
  id := c.Params("id")
  if id == "" {
    return c.Status(404).SendString("not found")
  }
  
  stpath := payload.StageGet(id)
  if stpath == "" {
    return c.Status(404).SendString("not found")
  }

  return c.SendFile(stpath)
}
