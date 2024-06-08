package agent

import (
	"fmt"

	"github.com/ngn13/ezcat/server/global"
	"github.com/ngn13/ezcat/server/util"
)

//            31     31           1         1       max 63
//           ==== =========   =========   =====   ==========
// request:  [id] [session] . [command] . [end] . [argument]

type Request struct {
  ID       string
  Session  string
  Command  byte 
  End      bool
  Argument string
}

var (
  REQ_REGISTER byte = 'R'
  REQ_CONTINUE byte = 'C'
  REQ_WORK     byte = 'W'
  REQ_DONE     byte = 'D'
  REQ_FAIL     byte = 'F'
)

func RequestParse(sections []string) (*Request, error) {
  if len(sections) != SECTION_COUNT {
    return nil, fmt.Errorf("not a valid request (section count: %d != %d)", len(sections), SECTION_COUNT)
  }

  if len(sections[0]) != FIRST_SECTION_LEN {
    return nil, fmt.Errorf("not a valid request (first section length: %d != %d)", len(sections[0]), FIRST_SECTION_LEN)
  }

  if len(sections[1]) != SECOND_SECTION_LEN {
    return nil, fmt.Errorf("not a valid request (second section length: %d != %d)", len(sections[1]), SECOND_SECTION_LEN)
  }

  if len(sections[2]) != THIRD_SECTION_LEN {
    return nil, fmt.Errorf("not a valid request (third section length: %d != %d)", len(sections[2]), THIRD_SECTION_LEN)
  }

  if len(sections[3]) > ARGUMENT_LEN {
    return nil, fmt.Errorf("not a valid request (argument is too large: %d > %d)", len(sections[3]), ARGUMENT_LEN)
  }

  var req Request = Request{
    ID:       util.Rot13(sections[0][:global.ID_LEN]),
    Session:  util.Rot13(sections[0][FIRST_SECTION_LEN-global.SESSION_LEN:]),
    Command:  util.Rot13Byte(sections[1][0]),
    End:      util.Rot13(sections[2])=="1",
    Argument: util.Rot13(sections[3]),
  }

  return &req, nil
}
