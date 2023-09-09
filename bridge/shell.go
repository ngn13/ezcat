package bridge

type Shell struct {
  Connected bool
  Hostname  string
  Lastcon   string
  Result    string
  Lhost     string
  Cmd       string
  ID        string
  Port      int
}

var SHELLS []Shell = []Shell{}

func CheckID(id string) int {
  for i, s := range SHELLS {
    if s.ID == id {
      return i
    }
  }

  return -1
}

func GetAvPort() int {
  start := 1234

  for _, s := range SHELLS {
    if s.Port == start {
      start -= 1 
    }
  }

  return start
}
