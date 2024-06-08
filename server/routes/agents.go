package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/agent"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

func GET_agents(c *fiber.Ctx) error {
  var res []agent.Data
  agent.Clean()

  for i := range agent.List {
    cur := &agent.List[i]

    if !cur.Active {
      continue
    }

    if cur.Username == "" || cur.Hostname == "" {
      log.Debug("Not listing %s (info not received yet)", cur.ID)
      continue
    }

    res = append(res, cur.Data())
  }

  return c.JSON(&fiber.Map{
    "list": res,
  })
}

func PUT_run(c *fiber.Ctx) error {
  agent.Clean()

  var data map[string]string
  if err := c.BodyParser(&data); err != nil {
    return util.ErrorCode(c, 400)
  }

  if data["address"] == "" || data["id"] == "" {
    return util.ErrorCode(c, 400)
  }

  ip, port, err := util.ParseAddr(data["address"])
  if err != nil {
    return util.ErrorCode(c, 400)
  }

  a := agent.Get(data["id"])
  if a == nil {
    return util.ErrorCode(c, 404)
  }

  if !a.Active {
    return util.Error(c, "Agent is not active")
  }

  work := a.AddWork(agent.CMD_RUN, fmt.Sprintf("%s:%d", ip, port), nil)
  return c.JSON(&fiber.Map{
    "job": work.Job.ID,
  })
}

func KillCallack(w *agent.Work) {
  a := agent.Get(w.Session)
  agent.DefaultCallack(w)
  a.Deactivate()
}

func GET_kill(c *fiber.Ctx) error {
  agent.Clean()
  id := c.Query("id")

  if id == "" {
    return util.ErrorCode(c, 400)
  }

  a := agent.Get(id)
  if a == nil {
    return util.ErrorCode(c, 404)
  }

  if !a.Active {
    return util.Error(c, "Agent is not active")
  }

  work := a.AddWork(agent.CMD_KILL, "plz", KillCallack)
  return c.JSON(&fiber.Map{
    "job": work.Job.ID,
  })
}
