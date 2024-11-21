package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/builder"
)

func GET_stage(c *fiber.Ctx) error {
	var (
		id         string
		stage_path string
	)

	build := c.Locals("builder").(*builder.Struct)

	if id = c.Params("id"); id == "" {
		return c.Status(404).SendString("not found")
	}

	if stage_path = build.GetStage(id); stage_path == "" {
		return c.Status(404).SendString("not found")
	}

	return c.SendFile(stage_path)
}
