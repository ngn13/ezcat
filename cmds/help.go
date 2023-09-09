package cmds

import (
	"fmt"
	"strings"

	"github.com/ngn13/ezcat/term"
)

func Help(args []string){
  if len(args) != 0 {
    term.Error("Command does not take any arguments")
    return
  }

  term.Newline()
  out := "ORANGECommand     DescriptionRESET\n"
  out += "ORANGE=======     ===========RESET\n"

  for _, c := range CMDS {
    out += fmt.Sprintf("%s%s%s\n", c.Name, strings.Repeat(" ", 12-len(c.Name)), c.Desc)
  }
    
  out += fmt.Sprintf("exit        Quit the program\n")
  term.Print(out)
}
