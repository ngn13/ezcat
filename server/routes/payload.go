package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/builder"
	"github.com/ngn13/ezcat/server/config"
	"github.com/ngn13/ezcat/server/util"
)

func GET_payloads(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{
		"list": c.Locals("builder").(*builder.Struct).Payloads,
	})
}

func GET_address(c *fiber.Ctx) error {
	conf := c.Locals("config").(*config.Struct)
	port := c.Query("port")

	if port == "" {
		return c.JSON(&fiber.Map{
			"address": fmt.Sprintf("%s:%d", util.GetIP(), conf.HTTP_Port),
		})
	}

	return c.JSON(&fiber.Map{
		"address": fmt.Sprintf("%s:%s", util.GetIP(), port),
	})
}

func PUT_build(c *fiber.Ctx) (err error) {
	var (
		data    map[string]string
		target  *builder.Target
		payload *builder.Payload
		res     string
	)

	build := c.Locals("builder").(*builder.Struct)

	if err := c.BodyParser(&data); err != nil {
		return util.ErrorCode(c, 400)
	}

	if data["address"] == "" || data["type"] == "" || data["os"] == "" {
		return util.ErrorCode(c, 400)
	}

	if target = builder.TargetByCode(data["os"]); target == nil {
		return util.Error(c, "Unknown OS/Arch")
	}

	if payload = build.GetPayload(data["type"]); payload == nil {
		return util.Error(c, "Unknown payload type")
	}

	if res, err = build.Create(payload, target, data["address"]); err != nil {
		return util.Error(c, fmt.Sprintf("Build failed: %s", err.Error()))
	}

	return c.JSON(&fiber.Map{
		"payload": string(res),
	})
}
