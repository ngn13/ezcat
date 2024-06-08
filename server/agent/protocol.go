package agent

import (
	"net"

	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/payload"
)

var (
  SECTION_COUNT      int = 4
  FIRST_SECTION_LEN  int = 62
  SECOND_SECTION_LEN int = 1
  THIRD_SECTION_LEN  int = 1
  ARGUMENT_LEN       int = 63
)

func ProtocolHandle(client *net.UDPAddr, data []string) (string, error) {
  req, err := RequestParse(data)
  if err != nil {
    return "", err
  }
  
  if payload.StageGet(req.ID) == "" {
    return ResponseDump(&Response{
      Command: RES_FAIL,
      Argument: "Invalid ID",
    })
  }

  agent := Get(req.Session)

  if agent == nil && req.Session == req.ID && req.Command == REQ_REGISTER {
    agent = New(req.ID)
    log.Debug("New agent with ID: %s and Session: %s", agent.ID, agent.Session)
    agent.IP = client.String()

    return ResponseDump(&Response{
      Command: RES_OK,
      Argument: agent.Session,
    })
  }

  if agent == nil {
    return ResponseDump(&Response{
      Command: RES_FAIL,
      Argument: "Agent not registered",
    })
  }

  err = agent.NewRequest(req)

  if err != nil {
    log.Debug("Failed to add new request to agent: %s", err.Error())
    return ResponseDump(&Response{
      Command: RES_FAIL,
      Argument: err.Error(),
    })
  }

  if agent.IsEnd() {
    return ProtocolHandleAgent(agent)
  }

  return "", nil 
}

func ProtocolHandleAgent(agent *Agent) (string, error) {
  req := agent.Request

  switch req.Command {
  case REQ_REGISTER:
    return ResponseDump(&Response{
      Command: RES_FAIL,
      Argument: "Already registered",
    })

  case REQ_WORK:
    work := agent.GetWork()
    if nil == work {
      return ResponseDump(&Response{
        Command: RES_NOTNOW,
        Argument: "Ask again",
      })
    }

    return ResponseDump(&Response{
      Command: work.Command,
      Argument: work.Argument,
    })

  case REQ_DONE:
    work := agent.GetWork()
    if work == nil {
      log.Debug("%s says it completed a work that does not exists, maybe the packet was delayed?", agent.ID)

      return ResponseDump(&Response{
        Command: RES_FAIL,
        Argument: "What are you talking about?",
      })
    }

    work.Response = agent.Request.Argument
    work.Success  = true
    work.Waiting  = false

    agent.HandleWork()
    agent.DelWork(work)

    return ResponseDump(&Response{
      Command: RES_OK,
      Argument: "Good job",
    })

  case REQ_FAIL:
    work := agent.GetWork()
    if work == nil {
      log.Debug("%s says it failed a work that does not exists, maybe the packet was delayed?", agent.ID)

      return ResponseDump(&Response{
        Command: RES_FAIL,
        Argument: "What are you talking about?",
      })
    }

    work.Response = agent.Request.Argument
    work.Success  = false 
    work.Waiting  = false

    agent.HandleWork()
    agent.DelWork(work)

    return ResponseDump(&Response{
      Command: RES_OK,
      Argument: "Damn thats unfortunate",
    })
  }

  return "", nil
} 
