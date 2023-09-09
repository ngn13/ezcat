package cmds

import (
	"fmt"

	"github.com/ngn13/ezcat/bridge"
	"github.com/ngn13/ezcat/term"
)

func List(args []string){
  if len(args) != 0 {
    term.Error("Command does not take any arguments")
    return
  }

  results := []bridge.Shell{}

  for _, s := range bridge.SHELLS {
    if s.Connected {
      results = append(results, s)
    }
  }

  if len(results) == 0 {
    term.Error("No active shells :(")
    return
  }

  term.Newline()
  out := "ORANGEIDs       Username/Hostname    Last ConnectionRESET\n"
  out += "ORANGE======    =================    ===============RESET\n"

  for _, s := range results {
    out += fmt.Sprintf("%s    %s          %s\n", s.ID, s.Hostname, s.Lastcon)
  }

  term.Print(out)
}
