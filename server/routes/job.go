package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/agent"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

func GET_job(c *fiber.Ctx) error {
	var (
		job *agent.Job
		id  uint16
		err error
	)

	list := c.Locals("agents").(*agent.List)
	list.Update()

	if id, err = util.ToUint16(c.Query("id")); err != nil {
		log.Debg("failed to parse job ID: %s", err.Error())
		return util.Error(c, "Invalid job ID")
	}

	if job = list.GetJob(id); job == nil {
		return util.Error(c, "Job not found")
	}

	res := fiber.Map{
		"id":      job.ID,
		"waiting": job.Waiting,
		"success": job.Success,
	}

	if job.Response != nil {
		res["message"] = string(job.Response)
	}

	return c.JSON(&res)
}

func DEL_job(c *fiber.Ctx) error {
	var (
		id  uint16
		err error
	)

	list := c.Locals("agents").(*agent.List)
	list.Update()

	if id, err = util.ToUint16(c.Query("id")); err != nil {
		log.Debg("failed to parse job ID: %s", err.Error())
		return util.Error(c, "Invalid job ID")
	}

	list.DelJob(id)
	return c.JSON(nil)
}
