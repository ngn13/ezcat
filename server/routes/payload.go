package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/global"
	"github.com/ngn13/ezcat/server/jobs"
	"github.com/ngn13/ezcat/server/payload"
	"github.com/ngn13/ezcat/server/util"
)

func GET_payloads(c *fiber.Ctx) error {
  return c.JSON(&fiber.Map{
    "list": payload.List,
  })
}

func GET_address(c *fiber.Ctx) error {
  port := c.Query("port")
  if port == "" {
    return c.JSON(&fiber.Map{
      "address": fmt.Sprintf("%s:%d", util.GetIP(), global.CONFIG_HTTPPORT),
    })
  }

  return c.JSON(&fiber.Map{
    "address": fmt.Sprintf("%s:%s", util.GetIP(), port),
  })
}

func PUT_build(c *fiber.Ctx) error {
  var data map[string]string

  if err := c.BodyParser(&data); err != nil {
    return util.ErrorCode(c, 400)
  }

  if data["address"] == "" || data["type"] == "" || data["os"] == "" {
    return util.ErrorCode(c, 400)
  }

  target := payload.TargetByCode(data["os"])
  if target == nil {
    return util.Error(c, "Unknown OS/Arch")
  }

  pyld := payload.Get(data["type"])
  if pyld == nil {
    return util.Error(c, "Unknown payload type")
  }

  job := jobs.Add("Building the payload...")

  go func(){
    res, err   := pyld.Build(target, data["address"])
    job.Active  = false
    util.CleanTemp()

    if err != nil {
      job.Message = fmt.Sprintf("Build failed: %s", err.Error())
      job.Success = false
      return
    }

    job.Message = res 
    job.Success = true 
  }()

  return c.JSON(&fiber.Map{
    "job": job.ID,
  })
}
