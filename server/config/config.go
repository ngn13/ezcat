package config

import (
	"os"
	"strconv"

	"github.com/ngn13/ezcat/server/global"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

func Load() {
  var err error

  global.CONFIG_HTTPPORT, err = strconv.Atoi(os.Getenv("HTTP_PORT"))
  if err != nil || !util.IsValidPort(global.CONFIG_HTTPPORT) {
    global.CONFIG_HTTPPORT = 5566
  }

  global.CONFIG_AGENTPORT, err = strconv.Atoi(os.Getenv("AGENT_PORT"))
  if err != nil || !util.IsValidPort(global.CONFIG_AGENTPORT) {
    global.CONFIG_AGENTPORT = 1053
  }

  global.CONFIG_PASSWORD = os.Getenv("PASSWORD")
  if global.CONFIG_PASSWORD == "" {
    log.Warn("Using the default password (very insecure, change it using the PASSWORD environment variable)")
    global.CONFIG_PASSWORD = "ezcat"
  }

  global.CONFIG_STATICDIR = os.Getenv("STATIC_DIR")
  if global.CONFIG_STATICDIR == "" {
    log.Warn("STATICDIR is not set, only serving the API")
  }

  global.CONFIG_DATADIR = os.Getenv("DATA_DIR")
  if global.CONFIG_DATADIR == "" {
    global.CONFIG_DATADIR = "./data"
  }

  global.CONFIG_PAYLOADDIR = os.Getenv("PAYLOAD_DIR")
  if global.CONFIG_PAYLOADDIR == "" {
    global.CONFIG_PAYLOADDIR = "../payloads"
  }

  global.CONFIG_MEGAMIND = os.Getenv("DISABLE_MEGAMIND") != "1"
  global.CONFIG_DEBUG = os.Getenv("DEBUG") == "1"
}
