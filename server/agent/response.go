package agent

import (
	"fmt"

	"github.com/ngn13/ezcat/server/util"
)

//               1       max 63
//           ========= ==========
// response: [command] [argument]

type Response struct {
  Command  byte 
  Argument string
}

var (
  RES_NOTNOW   byte = 'N'
  RES_FAIL     byte = 'F'
  RES_OK       byte = 'K'
  
  RES_INFO byte = 'I'
  CMD_INFO = RES_INFO 

  RES_RUN byte = 'R'
  CMD_RUN = RES_RUN

  RES_KILL byte = 'D'
  CMD_KILL = RES_KILL
)

func ResponseDump(res *Response) (string, error) {
  if len(res.Argument) > ARGUMENT_LEN{
    return "", fmt.Errorf("not a valid response (argument is too large: %d > %d)", len(res.Argument), ARGUMENT_LEN)
  }

  return util.Rot13(string(res.Command)+res.Argument), nil
}
