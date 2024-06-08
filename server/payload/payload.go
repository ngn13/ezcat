package payload

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ngn13/ezcat/server/global"
	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

type Payload struct {
  Name    string   `json:"name"`
  Config  string   `json:"config"`
  Output  string   `json:"output"`
  OS      []Target `json:"os"`
}

var List []Payload = []Payload{
  // linux only
  {Name: "Bash",       OS: []Target{TARGET_LINUX},  Config: "src/main.sh",   Output: "dist/main.sh"},

  // windows only
  {Name: "Powershell", OS: []Target{TARGET_WINDOWS}, Config: "src/main.ps1",  Output: "dist/main.ps1"},

  // both
  {Name: "PHP",        OS: TARGET_ALL, Config: "src/main.php", Output: "dist/main.php"},
  {Name: "Python",     OS: TARGET_ALL, Config: "src/main.py",  Output: "dist/main.py"},
}

func (p *Payload) Build(target *Target, address string) (string, error) {
  id  := util.MakeRandom(global.ID_LEN)
  dst := path.Join("/tmp", "ezcat_payload_"+id)

  ret := p.SupportsTarget(target)
  if !ret {
    return "", fmt.Errorf("target is not supported by the payload")
  }

  addr, port, err := util.ParseAddr(address)
  if err != nil {
    return "", err
  }

  _, err = StageBuild(StageConfig{
    Target:    target,
    ID:        id,
    Address:   addr,
    Port:      global.CONFIG_AGENTPORT,
  })

  if err != nil {
    return "", err
  }

  err = util.CopySimple(dst, p.GetDir(target.CodeName()))
  if err != nil {
    return "", err
  }

  config_file := path.Join(dst, p.Config)
  config_data, err := os.ReadFile(config_file)
  if err != nil {
    return "", err
  }

  config_url := fmt.Sprintf("http://%s:%d/%s", addr, port, id)
  config_edited := strings.ReplaceAll(string(config_data), "#URL#", config_url)
  err = os.WriteFile(config_file, []byte(config_edited), os.ModePerm)
  if err != nil {
    return "", err
  }

  out, err := util.RunBuild(dst, target.CodeName())
  if err != nil {
    log.Err("Failed to run the build command for the payload:\n=====\n%s%s\n=====", out, err.Error())
    return "", fmt.Errorf("failed to run the build: %s", err.Error())
  }

  res, err := os.ReadFile(path.Join(dst, p.Output))
  if err != nil {
    return "", err
  }

  if res[len(res)-1] == '\n' {
    res = res[:len(res)-1]
  }

  return string(res), nil
}

func (p *Payload) Supports(oscode string) *Target {
  target := TargetByCode(oscode)
  if target == nil {
    return nil
  }

  for i := range p.OS {
    if p.OS[i].CodeName() == target.CodeName() {
      return &p.OS[i]
    }
  }

  return nil
}

func (p *Payload) SupportsTarget(target *Target) bool {
  for i := range p.OS {
    if p.OS[i].CodeName() == target.CodeName() {
      return true
    }
  }

  return false
}

func (p *Payload) GetDir(oscode string) string {
  if p.Supports(oscode) == nil {
    return ""
  }

  src := path.Join(global.CONFIG_PAYLOADDIR, strings.ToLower(p.Name), oscode)
  if util.Access(src) {
    return src
  }

  src = path.Join(global.CONFIG_PAYLOADDIR, strings.ToLower(p.Name))
  if util.Access(src) {
    return src
  }

  return ""
}

func Get(typ string) *Payload{
  for i := range List {
    if List[i].Name == typ {
      return &List[i]
    }
  }
  return nil
}
