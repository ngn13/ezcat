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

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/agent"
	"github.com/ngn13/ezcat/server/config"
	"github.com/ngn13/ezcat/server/global"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/payload"
	"github.com/ngn13/ezcat/server/routes"
	"github.com/ngn13/ezcat/server/util"
)

func main() {
  app := fiber.New(fiber.Config{
    AppName: "ezcat",
    DisableStartupMessage: true,
  })

  // load config from the env
  config.Load()
  if global.CONFIG_STATICDIR != "" {
    app.Static("/", global.CONFIG_STATICDIR)
  }

  // agent server setup
  payload.StageLoad()
  var server agent.AgentServer
  go server.Start()

  // groups
  api  := app.Group("/api")
  user := api.Group("/user")

  // middlewares
  api.Use("*", util.CORS)
  user.Use("*", routes.ALL_auth)

  // user API routes
  user.Get("/logout",      routes.GET_logout)
  user.Get("/agent/list",  routes.GET_agents)
  user.Get("/agent/kill",  routes.GET_kill)
  user.Put("/agent/run",   routes.PUT_run)

  user.Get("/payload/list",  routes.GET_payloads)
  user.Get("/payload/addr",  routes.GET_address)
  user.Put("/payload/build", routes.PUT_build)

  user.Get("/job/get",    routes.GET_job)
  user.Delete("/job/del", routes.DEL_job)

  // other API routes (no auth needed)
  api.Get("/info",  routes.GET_info)
  api.Put("/login", routes.PUT_login)

  // payload routes
  app.Get("/:id",   routes.GET_stage)

  // start the api
  log.Info("Starting ezcat 🐱(v%s)", routes.VERSION)
  if global.CONFIG_STATICDIR != "" {
    log.Info("======> Visit http://127.0.0.1:%d <======", global.CONFIG_HTTPPORT)
  }
  log.Err(app.Listen(fmt.Sprintf(":%d", global.CONFIG_HTTPPORT)).Error())
}
