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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/ngn13/ezcat/log"
	"github.com/ngn13/ezcat/routes"
)

func main() {
  engine := html.New("./views", ".html")
  app := fiber.New(fiber.Config{
    Views: engine,
    AppName: "ezcat",
    DisableStartupMessage: true,
  })

  // static files
  app.Static("static", "./static")

  // admin routes
  app.Get("/", routes.GETLogin)
  app.Post("/", routes.POSTLogin)
  app.Use("/admin/*", routes.MIDAdmin)
  app.Get("/admin", routes.GETAdmin)
  app.Get("/admin/status", routes.GETStatus)
  app.Get("/admin/run", routes.GETRun)
  app.Post("/admin/run", routes.POSTRun)
  app.Get("/admin/clean", routes.GETClean)
  app.Get("/logout", routes.GETLogout)

  // shell routes
  app.Use("/shell/*", routes.MIDShell)
  app.Get("/shell/job", routes.GETJob)
  app.Post("/shell/result", routes.POSTRes)

  log.Info("Starting ezcat 🐱")
  log.Info("Visit http://127.0.0.1:5566")
  log.Err(app.Listen(":5566").Error())
}
