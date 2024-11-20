package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/config"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

var tokens []string = []string{}

func tokenOK(c *fiber.Ctx) bool {
	token := util.GetToken(c)
	for _, t := range tokens {
		if t == token {
			return true
		}
	}
	return false
}

func tokenNew() string {
	tokens = append(tokens, util.MakeRandom(31))
	return tokens[len(tokens)-1]
}

func tokenDel(token string) {
	for i, t := range tokens {
		if t == token {
			tokens = append(tokens[:i], tokens[i+1:]...)
			return
		}
	}
}

func Auth(c *fiber.Ctx) error {
	if tokenOK(c) {
		return c.Next()
	}

	return util.ErrorCode(c, 401)
}

func PUT_login(c *fiber.Ctx) error {
	var (
		data map[string]string
		conf *config.Struct
	)

	conf = c.Locals("config").(*config.Struct)

	if err := c.BodyParser(&data); err != nil {
		return util.ErrorCode(c, 400)
	}

	if data["password"] == "" {
		return util.ErrorCode(c, 400)
	}

	if conf.Password != data["password"] {
		return util.Error(c, "Invalid password")
	}

	log.Info("new login from %s", c.IP())

	return c.JSON(fiber.Map{
		"token": tokenNew(),
	})
}

func GET_logout(c *fiber.Ctx) error {
	tokenDel(util.GetToken(c))
	return c.JSON(fiber.Map{})
}
