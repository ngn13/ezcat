package log

import (
	"log"
	"os"

	"github.com/ngn13/ezcat/server/global"
)

var (
  Debug = func(f string, args...interface{}){ 
    if !global.CONFIG_DEBUG { 
      return 
    } 

    log.SetFlags(log.Ltime | log.Lshortfile)
    log.SetPrefix("\033[1m\033[36m[DEBG]\033[0m ")
    log.Printf(f, args...)
  }
  Warn  = log.New(os.Stdout, "\033[1m\033[33m[WARN]\033[0m ", log.Ltime | log.Lshortfile).Printf
  Info  = log.New(os.Stdout, "\033[1m\033[34m[INFO]\033[0m ", log.Ltime | log.Lshortfile).Printf
  Err   = log.New(os.Stderr, "\033[1m\033[31m[ERRO]\033[0m ", log.Ltime | log.Lshortfile).Printf
)
