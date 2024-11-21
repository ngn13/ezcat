package builder

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

type Payload struct {
	Name   string   `json:"name"`
	Config string   `json:"config"`
	Output string   `json:"output"`
	OS     []Target `json:"os"`
}

func (p *Payload) SupportsTarget(target *Target) bool {
	for i := range p.OS {
		if p.OS[i].CodeName() == target.CodeName() {
			return true
		}
	}

	return false
}

func (p *Payload) Supports(oscode string) *Target {
	var target *Target

	if target = TargetByCode(oscode); target == nil {
		return nil
	}

	if p.SupportsTarget(target) {
		return target
	}

	return nil
}

func (p *Struct) buildPayload(b *build) (string, error) {
	if !b.Payload.SupportsTarget(b.Target) {
		return "", fmt.Errorf("target is not supported by the payload")
	}

	payloaddir := path.Join(p.Config.PayloadDir, strings.ToLower(b.Payload.Name))
	builddir := path.Join("/tmp", "ezcat_payload_"+b.ID)

	if err := util.CopySimple(builddir, payloaddir); err != nil {
		return "", err
	}

	config_file := path.Join(builddir, b.Payload.Config)
	config_data, err := os.ReadFile(config_file)

	if err != nil {
		return "", err
	}

	config_url := fmt.Sprintf("http://%s:%d/%s", b.Host, b.Port, b.ID)
	config_edited := strings.ReplaceAll(string(config_data), "#URL#", config_url)

	if err = os.WriteFile(config_file, []byte(config_edited), os.ModePerm); err != nil {
		return "", err
	}

	if out, err := util.RunBuild(builddir, b.Target.CodeName(), nil); err != nil {
		log.Fail("failed to run the build command for the payload:\n=====\n%s%s\n=====", out, err.Error())
		return "", fmt.Errorf("failed to run the build: %s", err.Error())
	}

	res, err := os.ReadFile(path.Join(builddir, b.Payload.Output))

	if err != nil {
		return "", err
	}

	if res[len(res)-1] == '\n' {
		res = res[:len(res)-1]
	}

	return string(res), nil
}
