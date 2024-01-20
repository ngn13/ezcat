package routes

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/log"
	"github.com/ngn13/ezcat/shell"
	"github.com/ngn13/ezcat/util"
)

var Token string = util.MakeRandom(42)
var Status string = util.MakeRandom(13)

func UpdateStatus() {
  Status = util.MakeRandom(13)
}

func MIDShell(c *fiber.Ctx) error {
  token := c.Query("t")
  if token != Token {
    c.Status(403)
    return util.SendEnc(c, "echo QUIT")
  }

  return c.Next()
}

func ParseData(i int, d []byte) {
  l := strings.Split(string(d), ":")
  if len(l) != 2 {
    return
  }

  shell.List[i].Host = l[0]
  shell.List[i].User = l[1]
  shell.List[i].Success   = true

  log.Info("Got info for %s => User: %s, Host: %s", 
      shell.List[i].UID, shell.List[i].Host, shell.List[i].User)
  UpdateStatus()
}

func POSTRes(c *fiber.Ctx) error {
  uid  := c.Query("u")
  code := c.Query("c")

  if len(uid) != 7 || code == "" {
    return util.SendEnc(c, shell.S_QUIT)
  }

  indx := shell.Get(util.Reverse(uid)) 
  if indx == -1 {
    return util.SendEnc(c, shell.List[indx].Script)
  }

  shell.List[indx].Success = true 
  if code != "0" {
    shell.List[indx].Success = false
  }

  switch shell.List[indx].Script {
  case shell.S_GET_INFO:
    ParseData(indx, c.BodyRaw())
  }

  if shell.Update() {
    UpdateStatus()
  }
  shell.List[indx].Script = shell.S_PASS 
  return c.SendString("")
}

func GETJob(c *fiber.Ctx) error {
  uid := c.Query("u")
  if len(uid) != 7 {
    return util.SendEnc(c, shell.S_QUIT)
  }

  indx := shell.Get(util.Reverse(uid)) 
  if indx != -1 {
    shell.LastCon(indx)
    return util.SendEnc(
      c, shell.List[indx].Script)
  }
    
  news := shell.Shell{
    UID: util.Reverse(uid), 
    Script: shell.S_GET_INFO,
    IP: c.IP(),
    Host: "unknown",
    User: "unknown",
    Last: time.Now(),
  }

  shell.List = append(shell.List, news)
  log.Info("Got new shell %s => IP: %s", 
    news.UID, news.IP)
  UpdateStatus()

  return util.SendEnc(c, news.Script)
}
