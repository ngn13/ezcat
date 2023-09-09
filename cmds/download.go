package cmds

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ngn13/ezcat/bridge"
	"github.com/ngn13/ezcat/term"
)

func Download(args []string){
  if len(args) != 2 {
    term.Error("Please provide a shell ID and a file name")
    return
  }

  id := args[0]
  file := args[1]

  indx := bridge.CheckID(id)
  if indx == -1 || !bridge.SHELLS[indx].Connected {
    term.Error("Shell not found")
    return
  }
  
  bridge.SHELLS[indx].Cmd = fmt.Sprintf("cat %s", file)
  term.Success("Sent the download command")
  term.Info("Waiting for response")
  res := bridge.SHELLS[indx].Result 

  for {
    if res != bridge.SHELLS[indx].Result {
      res = bridge.SHELLS[indx].Result
      break
    }
    time.Sleep(1)
  }

  if strings.Contains(strings.ToLower(res), "no such file"){
    term.Error("File not found on the remote machine")
    return
  }
  
  term.Info("Received the file!")
  err := os.WriteFile(file, []byte(res), 0644)
  if err != nil {
    term.Error("Error writing file: %s", err)
    return
  }

  term.Success("File downloaded successfuly")
}
