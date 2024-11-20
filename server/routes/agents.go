package routes

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/agent"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

func GET_agents(c *fiber.Ctx) error {
	list := c.Locals("agents").(*agent.List)
	list.Update()

	return c.JSON(&fiber.Map{
		"list": list,
	})
}

func PUT_run(c *fiber.Ctx) error {
	var (
		data map[string]string

		list *agent.List
		ag   *agent.Agent

		id   uint64
		ip   string
		port uint16
		err  error
	)

	list = c.Locals("agents").(*agent.List)
	list.Update()

	if err = c.BodyParser(&data); err != nil {
		log.Debg("failed to parse the body: %s", err.Error())
		return util.ErrorCode(c, 400)
	}

	if data["address"] == "" || data["id"] == "" {
		return util.ErrorCode(c, 400)
	}

	if ip, port, err = util.ParseAddr(data["address"]); err != nil {
		log.Debg("failed to parse the address: %s", err.Error())
		return util.ErrorCode(c, 400)
	}

	if id, err = strconv.ParseUint(data["session"], 10, 32); err != nil {
		log.Debg("failed to parse ID: %s", err.Error())
		return util.ErrorCode(c, 400)
	}

	if ag = list.Find(uint32(id)); ag == nil {
		return util.ErrorCode(c, 404)
	}

	if !ag.Conneceted {
		return util.Error(c, "agent is not active")
	}

	job := ag.AddJob(agent.CMD_RUN, fmt.Sprintf("%s:%d", ip, port), nil)

	return c.JSON(&fiber.Map{
		"job": job.ID,
	})
}

func KillCallack(j *agent.Job) {
	j.Agent.ShouldKill = true
}

func GET_kill(c *fiber.Ctx) error {
	var (
		ag  *agent.Agent
		id  uint64
		err error
	)

	list := c.Locals("agents").(*agent.List)
	list.Update()

	if id, err = strconv.ParseUint(c.Query("session"), 10, 32); err != nil {
		log.Debg("failed to parse ID: %s", err.Error())
		return util.ErrorCode(c, 400)
	}

	if ag = list.Find(uint32(id)); ag == nil {
		return util.ErrorCode(c, 404)
	}

	if !ag.Conneceted {
		return util.Error(c, "agent is not active")
	}

	job := ag.AddJob(agent.CMD_KILL, "plz", KillCallack)

	return c.JSON(&fiber.Map{
		"job": job.ID,
	})
}
