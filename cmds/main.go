package cmds

type Cmd struct {
  Name  string
  Desc  string
  Func  func([]string)
}

var CMDS []Cmd = []Cmd{}

func LoadCmds(){
  CMDS = append(CMDS, Cmd{
    Name: "help",
    Desc: "List all the commands",
    Func: Help,
  })

  CMDS = append(CMDS, Cmd{
    Name: "list",
    Desc: "List all shells",
    Func: List,
  })

  CMDS = append(CMDS, Cmd{
    Name: "gen",
    Desc: "Generate a shell payload",
    Func: Gen,
  })

  CMDS = append(CMDS, Cmd{
    Name: "shell",
    Desc: "Hop on a shell",
    Func: Shell,
  })

  CMDS = append(CMDS, Cmd{
    Name: "upload",
    Desc: "Upload a file",
    Func: Upload,
  })

  CMDS = append(CMDS, Cmd{
    Name: "download",
    Desc: "Download a file",
    Func: Download,
  })
}
