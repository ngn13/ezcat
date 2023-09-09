package term

import (
	"fmt"
	"strings"
)

var RESET  = "\033[0m"
var RED    = "\033[31m"
var GREEN  = "\033[32m"
var BLUE   = "\033[34m"
var PURPLE = "\033[35m"
var CYAN   = "\033[36m"
var GRAY   = "\033[37m"
var WHITE  = "\033[97m"
var ORANGE = "\033[38;5;208m"

func Newline(){
  fmt.Printf("\n")
}

func Print(format string, v ...interface{}) {
  redd := strings.ReplaceAll(format, "RED", RED)
  blued := strings.ReplaceAll(redd, "BLUE", BLUE)
  greend := strings.ReplaceAll(blued, "GREEN", GREEN)
  purpled := strings.ReplaceAll(greend, "PURPLE", PURPLE)
  cyand := strings.ReplaceAll(purpled, "CYAN", CYAN)
  grayd := strings.ReplaceAll(cyand, "GRAY", GRAY)
  oranged := strings.ReplaceAll(grayd, "ORANGE", ORANGE)
  resetd := strings.ReplaceAll(oranged, "RESET", RESET)

  fmt.Printf(GREEN+resetd+RESET+"\n", v...)
}

func Success(format string, v ...interface{}) {
  fmt.Printf(GREEN+"[+] "+RESET+format+"\n", v...)
}

func Error(format string, v ...interface{}) {
  fmt.Printf(ORANGE+"[-] "+RESET+format+"\n", v...)
}

func Info(format string, v ...interface{}) {
  fmt.Printf(BLUE+"[*] "+RESET+format+"\n", v...)
}

func Cyan(format string, v ...interface{}) {
  fmt.Printf(CYAN+format+RESET+"\n", v...)
}
func Gray(format string, v ...interface{}) {
  fmt.Printf(GRAY+format+RESET+"\n", v...)
}
