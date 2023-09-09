package shell

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/ngn13/ezcat/bridge"
	"github.com/ngn13/ezcat/term"
)

func StartListener() error{
  port := bridge.GetAvPort()
  term.Print("BLUE[*]RESET Starting netcat listener on port ORANGE%dRESET", port)
  cmd := exec.Command("nc", "-lnvp", strconv.Itoa(port))

  cmd.Stdout = os.Stdout
  cmd.Stdin = os.Stdin
  cmd.Stderr = os.Stderr
  cmd.Run()

  /*
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return errors.New("Error getting stdout pipe: %s"+err.Error())
  }

  stdin, err := cmd.StdinPipe()
  if err != nil {
    return errors.New("Error getting stdin pipe: %s"+err.Error())
  }

  err = cmd.Start()
  if err != nil {
    return errors.New("Error running netcat listener: %s"+err.Error())
  }

  scanner := bufio.NewScanner(stdout)
  for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line) 
    if strings.Contains(line, "Connection from") {
      stdin.Write([]byte("python3 -c 'import pty;pty.spawn(\"/bin/bash\")"))
      stdin.Write([]byte("export TERM=xterm"))
      break
    }
  }*/

  cmd.Wait()
  term.Info("Exit with the code %d", cmd.ProcessState.ExitCode())
  return nil
}
