/*
 *  ezcat | easy reverse shell handler
 *  written by ngn (https://ngn.tf) (2024)
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/agent"
	"github.com/ngn13/ezcat/server/builder"
	"github.com/ngn13/ezcat/server/c2"
	"github.com/ngn13/ezcat/server/config"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/routes"
)

func main() {
	var (
		agents agent.List
		build  *builder.Struct
		conf   *config.Struct
		app    *fiber.App
		srv    *c2.Server
		err    error
	)

	// load config from the env
	if conf, err = config.New(); err != nil {
		log.Fail("failed to load the configuration: %s", err.Error())
		os.Exit(1)
	}

	// load the payload/stage builder
	if build, err = builder.New(conf); err != nil {
		log.Fail("failed to create a new builder: %s", err.Error())
		os.Exit(1)
	}

	// create a C2 server
	srv = c2.New(build, &agents)

	// create a web server
	app = fiber.New(fiber.Config{
		AppName:               "ezcat",
		ServerHeader:          "",
		DisableStartupMessage: true,
	})

	if conf.StaticDir != "" {
		app.Static("/", conf.StaticDir)
	}

	if !conf.Debug {
		log.Debg = func(format string, v ...any) {}
	}

	// groups
	api := app.Group("/api")
	user := api.Group("/user")

	// middlewares
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("agents", &agents)
		c.Locals("builder", build)
		c.Locals("config", conf)
		return c.Next()
	})

	// auth and CORS middlewars
	api.Use(routes.CORS)
	user.Use(routes.Auth)

	// user API routes
	user.Get("/logout", routes.GET_logout)
	user.Get("/agent/list", routes.GET_agents)
	user.Get("/agent/kill", routes.GET_kill)
	user.Put("/agent/run", routes.PUT_run)

	user.Get("/payload/list", routes.GET_payloads)
	user.Get("/payload/addr", routes.GET_address)
	user.Put("/payload/build", routes.PUT_build)

	user.Get("/job/get", routes.GET_job)
	user.Delete("/job/del", routes.DEL_job)

	// other API routes (no auth needed)
	api.Get("/info", routes.GET_info)
	api.Put("/login", routes.PUT_login)

	// payload routes
	app.Get("/:id", routes.GET_stage)

	// other routes
	app.All("*", func(c *fiber.Ctx) error {
		return c.Redirect("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	})

	// start the API
	log.Info("starting ezcat 🐱(v%s)", conf.Version)

	log.Debg("starting the C2 server on port %d", conf.C2_Port)

	if err = srv.Listen(fmt.Sprintf("0.0.0.0:%d", conf.C2_Port)); err != nil {
		log.Fail("failed to start the C2 server")
		os.Exit(1)
	}

	if conf.StaticDir != "" {
		log.Info("======> visit http://127.0.0.1:%d <======", conf.HTTP_Port)
	}

	log.Debg("starting the HTTP server on port %d", conf.HTTP_Port)

	if err = app.Listen(fmt.Sprintf("0.0.0.0:%d", conf.HTTP_Port)); err != nil {
		log.Fail("failed to start the web server")
		os.Exit(1)
	}
}
