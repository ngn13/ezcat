package routes 

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/util"
)

var tokens []string = []string{}

func TokenOK(c *fiber.Ctx) bool {
  token := util.GetToken(c) 
  for _, t := range tokens {
    if t == token {
      return true
    }
  }
  return false
}

func TokenNew() string {
  tokens = append(tokens, util.MakeRandom(31))
  return tokens[len(tokens)-1]
}

func TokenDel(token string) {
  for i, t := range tokens {
    if t == token {
      tokens = append(tokens[:i], tokens[i+1:]...)
      return
    }
  }
}

func ALL_auth(c *fiber.Ctx) error {
  if TokenOK(c) {
    return c.Next()
  }

  return util.ErrorCode(c, 401)
}
