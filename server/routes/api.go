package routes

import "github.com/gofiber/fiber/v2"

func CORS(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers",
		"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Set("Access-Control-Allow-Methods", "OPTIONS, PUT, DELETE, GET")

	if c.Method() == "OPTIONS" {
		return c.SendString("")
	}

	return c.Next()
}
