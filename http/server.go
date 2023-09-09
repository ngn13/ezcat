package http

import (
	"io"
	"net/http"
	"time"

	"github.com/ngn13/ezcat/bridge"
)

func SendCmd(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    io.WriteString(w, "Bad request")
    return 
  }

  params := r.URL.Query()
  shell_id := params.Get("id")
  indx := bridge.CheckID(shell_id)
  
  if indx == -1 {
    io.WriteString(w, "ID not found")
    return
  }

  bridge.SHELLS[indx].Lastcon = time.Now().Format("01/02/2006 15:04:05")

  body, err := io.ReadAll(r.Body)
  if err != nil {
    io.WriteString(w, "Bad request")
    return 
  }
 
  if !bridge.SHELLS[indx].Connected {
    res := "Unknown"
    if len(body) != 0 {
      res = string(body)
    }
    //term.Cyan("\r\rNew connection from %s!", r.RemoteAddr)
    bridge.SHELLS[indx].Connected = true
    bridge.SHELLS[indx].Hostname = res
  }else {
    if len(body) != 0 {
      bridge.SHELLS[indx].Result = string(body)
    }
  }

  io.WriteString(w, bridge.SHELLS[indx].Cmd)
  bridge.SHELLS[indx].Cmd = ""
  return
} 

func Listen() error{
  addr := bridge.GetAddr()
  http.HandleFunc("/", SendCmd)
  err := http.ListenAndServe(addr, nil)
  return err
}
