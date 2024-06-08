package payload

import (
	"fmt"
	"strings"
)

type Target struct {
  Name string `json:"name"`
  Arch string `json:"arch"`
}

var (
  TARGET_LINUX   = Target{
    Name: "Linux",
    Arch: "amd64",
  }

  TARGET_WINDOWS = Target{
    Name: "Windows",
    Arch: "amd64",
  }

  TARGET_ALL     = []Target{
    TARGET_LINUX, TARGET_WINDOWS,
  }
)

func (t *Target) PrettyName() string{
  return fmt.Sprintf("%s (%s)", t.Name, t.Arch)
}

func (t *Target) CodeName() string {
  return fmt.Sprintf("%s_%s", 
    strings.ToLower(t.Name), strings.ToLower(t.Arch))
}

func TargetByCode(cn string) *Target {
  for i := range TARGET_ALL {
    if TARGET_ALL[i].CodeName() == cn {
      return &TARGET_ALL[i]
    }
  }
  return nil
}
