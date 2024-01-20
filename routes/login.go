package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/log"
	"github.com/ngn13/ezcat/util"
)

var Tokens []string

func AddToken() string {
  token := util.MakeRandom(42) 
  Tokens = append(Tokens, string(token))
  return string(token)
}

func RemoveToken(token string) {
  for i := range Tokens {
    if Tokens[i] != token {
      continue 
    }
    Tokens = append(Tokens[:i], Tokens[i+1:]...)
    break
  }
}

func CheckToken(token string) bool {
  for _, t := range Tokens {
    if t == token {
      return true
    }
  }

  return false
}

func POSTLogin(c *fiber.Ctx) error {
  body := struct {
    Pass string `form:"password"`
  }{}

  if err := c.BodyParser(&body); err != nil {
    return util.RenderErr(c, 500, err)
  }

  pass := os.Getenv("PASSWORD")
  if pass == "" {
    pass = "ezcat"
  }

  if pass == body.Pass {
    c.Cookie(&fiber.Cookie{
      Name: "token",
      Value:  AddToken(),
    })
    log.Info("New login from %s", c.IP())
    return c.Redirect("/admin")
  }
   
  log.Warn("Bad login from %s", c.IP())
  return c.Render("login", fiber.Map{
    "failed": true,
  })
}

func GETLogin(c *fiber.Ctx) error {
  token := c.Cookies("token")

  if !CheckToken(token) {
    return c.Render("login", fiber.Map{})
  }

  return c.Redirect("/admin")
}

func GETLogout(c *fiber.Ctx) error {
  token := c.Cookies("token")

  if !CheckToken(token) {
    return c.Redirect("/")
  }

  RemoveToken(token)
  c.ClearCookie()

  log.Info("Logout from %s", c.IP())
  return c.Redirect("/")
}
