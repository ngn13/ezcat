package log

import (
	"log"
	"os"
)

var (
  Warn = log.New(os.Stdout, "\033[1m\033[33m[WARN]\033[0m ", log.Ltime).Printf
  Info = log.New(os.Stdout, "\033[1m\033[34m[INFO]\033[0m ", log.Ltime).Printf
  Err  = log.New(os.Stderr, "\033[1m\033[31m[ERRO]\033[0m ", log.Ltime).Printf
)
