package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/agent"
	"github.com/ngn13/ezcat/server/jobs"
	"github.com/ngn13/ezcat/server/util"
)

func GET_job(c *fiber.Ctx) error {
  agent.Clean()

  id := c.Query("id")
  if id == "" {
    return util.Error(c, "Job ID is not specified")
  }

  job := jobs.Get(id)
  if job == nil {
    return util.Error(c, "Job not found")
  }

  return c.JSON(job)
}

func DEL_job(c *fiber.Ctx) error {
  id := c.Query("id")
  if id == "" {
    return util.Error(c, "Job ID is not specified")
  }

  jobs.Del(id)
  return c.JSON(&fiber.Map{})
}
