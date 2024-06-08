package global

// these defaults are reset in config/config.go
var (
  CONFIG_HTTPPORT   int    = 5566
  CONFIG_AGENTPORT  int    = 1053
  CONFIG_PASSWORD   string = "ezcat"
  CONFIG_STATICDIR  string = ""
  CONFIG_DATADIR    string = "./data"
  CONFIG_PAYLOADDIR string = "../payloads"
  CONFIG_DEBUG      bool   = false
  CONFIG_MEGAMIND   bool   = true
)
