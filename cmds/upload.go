package cmds

import (
	"fmt"
	"os"
	"strings"

	"github.com/ngn13/ezcat/bridge"
	"github.com/ngn13/ezcat/term"
)

func Upload(args []string){
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

  st, err := os.Stat(file)

  if os.IsNotExist(err){
    term.Error("File not found")
    return
  }else if err != nil {
    term.Error("Error accessing file: %s", err)
    return
  }

  if st.IsDir() {
    term.Error("Please provide a file, not a directory")
    return
  }

  raw, err := os.ReadFile(file)
  if err != nil {
    term.Error("Error reading the file: %s", err)
    return
  }

  // shitty filtering xd
  final := strings.ReplaceAll(string(raw), "\"", "\\\"")
  final = strings.ReplaceAll(string(final), "`", "\\`")
  final = strings.ReplaceAll(string(final), "$", "\\$")
  final = strings.ReplaceAll(string(final), "&", "\\&")

  bridge.SHELLS[indx].Cmd = fmt.Sprintf("echo \"%s\" > %s", final, file)
  term.Info(bridge.SHELLS[indx].Cmd)
  term.Success("Sent the upload command")
}
