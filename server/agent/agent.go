package agent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ngn13/ezcat/server/global"
	"github.com/ngn13/ezcat/server/jobs"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

type Work struct {
  Session  string
  Job      *jobs.Job

  Command  byte
  Argument string
  Response string
  Waiting  bool
  Success  bool
  Callback func(*Work)
}

type Agent struct {
  Username string
  Hostname string
  Kernel   string
  PID      int
  IP       string

  Lastcon  time.Time
  Active   bool

  Session  string // SESSION_LEN
  ID       string // ID_LEN
 
  Request  *Request
  Work     []Work
}

type Data struct {
  Username string `json:"username"`
  Hostname string `json:"hostname"`
  Kernel   string `json:"kernel"`
  PID      int    `json:"pid"`
  IP       string `json:"ip"`
  ID       string `json:"id"` // not really the "ID", its the session
}

var List []Agent

func InfoCallback(w *Work) {
  if !w.Success {
    log.Err("Failed to get info from the agent: %s", w.Response)
    jobs.Del(w.Job.ID)
    return
  }

  var err error

  data := strings.Split(w.Response, "@")
  agent := Get(w.Session)

  if len(data) != 4 {
    log.Debug("Bad info data from %s", agent.ID)
    log.Debug(w.Response)
    return
  }

  agent.Username  = data[0]
  agent.Hostname  = data[1]
  agent.Kernel    = data[2]

  agent.PID, err = strconv.Atoi(data[3])
  if err != nil {
    agent.PID = -1
  }

  jobs.Del(w.Job.ID)
}

func DefaultCallack(w *Work){
  w.Job.Message = w.Response
  w.Job.Active  = false
}

func New(id string) *Agent {
  agent := Agent{
    ID:       id,
    Session:  util.MakeRandom(global.SESSION_LEN),
    
    Active:   true,
    Lastcon:  time.Now(),

    Username: "",
    Hostname: "",
    Kernel:   "",
    IP:       "",

    Request:  nil,
  }

  agent.AddWork(RES_INFO, "plz", InfoCallback)
  List = append(List, agent)
  return &List[len(List)-1]
}

func Get(sess string) *Agent {
  for i := range List {
    if List[i].Session == sess {
      return &List[i]
    }
  }
  return nil
}

func Clean() {
  for i := range List {
    cur := &List[i]
    if time.Since(cur.Lastcon) > time.Second*time.Duration(global.SLEEP_MAX+5) {
      cur.Deactivate()
    }
  }
}

func (a *Agent) IsEnd() bool {
  if(a.Request == nil){
    return true 
  }
  return a.Request.End
}

func (a *Agent) NewRequest(r *Request) error {
  a.Lastcon = time.Now()
  a.Active  = true

  if a.IsEnd() && r.Command != REQ_CONTINUE {
    a.Request = r
    return nil
  }

  if r.Command != REQ_CONTINUE {
    return fmt.Errorf("invalid order")
  }

  a.Request.Argument += r.Argument
  a.Request.End = r.End
  return nil 
}

func (a *Agent) GetWork() *Work {
  for i := range a.Work {
    if a.Work[i].Waiting {
      return &a.Work[i]
    }
  }
  return nil
} 

func (a *Agent) AddWork(cmd byte, arg string, callback func(*Work)) *Work {
  job := jobs.Add("Waiting for response from the agent")
  if callback == nil {
    callback = DefaultCallack
  }

  a.Work = append(a.Work, Work{
    Session: a.Session,
    Job: job,

    Command:  cmd,
    Argument: arg,
    Response: "",
    Waiting:  true,
    Callback: callback,
  })

  return &a.Work[len(a.Work)-1]
}

func (a *Agent) HandleWork() {
  for i := range a.Work {
    if !a.Work[i].Waiting {
      a.Work[i].Job.Success = a.Work[i].Success
      a.Work[i].Callback(&a.Work[i])
    }
  }
}

func (a *Agent) DelWork(w *Work) {
  for i := range a.Work {
    if &a.Work[i] != w {
      continue
    }
    a.Work = append(a.Work[:i], a.Work[i+1:]...) 
    return
  }
}

func (a *Agent) Data() Data {
  return Data{
    Username: a.Username,
    Hostname: a.Hostname,
    Kernel:   a.Kernel,
    PID:      a.PID,
    IP:       a.IP,
    ID:       a.Session,
  }
}

func (a *Agent) Deactivate() {
  a.Active = false
  a.Request = nil
  
  for i := range a.Work {
    cur := a.Work[i]
    if !cur.Waiting {
      continue
    }

    cur.Waiting = false
    cur.Job.Active  = false
    cur.Job.Success = false
    cur.Job.Message = "Agent is not active"
  }

  a.Work = []Work{}
}
