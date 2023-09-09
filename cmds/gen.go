package cmds

import (
	"github.com/ngn13/ezcat/shell"
	"github.com/ngn13/ezcat/term"
)

func Gen(args []string){
  if len(args) != 1 {
    term.Error("Please provide an interface")
    return
  }

  out, err := shell.Gen(args[0])
  if err != nil {
    term.Error(err.Error())
    return
  }

  term.Info("Here you go buddy:")
  term.Gray(out)
}
