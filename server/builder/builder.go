package builder

import (
	"os"
	"path"

	"github.com/ngn13/ezcat/server/config"
	"github.com/ngn13/ezcat/server/util"
)

const ID_LEN = 32

type Struct struct {
	Config   *config.Struct
	Payloads []Payload
	Builds   []string
}

type build struct {
	ID      string
	Host    string
	Port    uint16
	Target  *Target
	Payload *Payload
}

func New(conf *config.Struct) (*Struct, error) {
	var (
		builder Struct
		entries []os.DirEntry
		err     error
	)

	builder = Struct{
		Config: conf,
		Payloads: []Payload{
			// linux only
			{Name: "Bash", OS: []Target{TARGET_LINUX}, Config: "src/main.sh", Output: "dist/main.sh"},

			// windows only
			{Name: "Powershell", OS: []Target{TARGET_WINDOWS}, Config: "src/main.ps1", Output: "dist/main.ps1"},

			// both
			{Name: "PHP", OS: TARGET_ALL, Config: "src/main.php", Output: "dist/main.php"},
			{Name: "Python", OS: TARGET_ALL, Config: "src/main.py", Output: "dist/main.py"},
		},
	}

	if entries, err = os.ReadDir(conf.DistDir); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	for _, e := range entries {
		builder.Builds = append(builder.Builds, e.Name())
	}

	return &builder, nil
}

func (s *Struct) GetPayload(name string) *Payload {
	for i := range s.Payloads {
		if s.Payloads[i].Name == name {
			return &s.Payloads[i]
		}
	}

	return nil
}

func (s *Struct) GetStage(id string) string {
	for _, i := range s.Builds {
		if i == id {
			return ""
		}
	}

	return path.Join(s.Config.DistDir, id)
}

func (s *Struct) Create(payload *Payload, target *Target, address string) (string, error) {
	var (
		res  string
		host string
		port uint16
		err  error
	)

	// parse the address
	if host, port, err = util.ParseAddr(address); err != nil {
		return "", err
	}

	// create the build configuration
	build := build{
		ID:      util.MakeRandom(ID_LEN),
		Host:    host,
		Port:    port,
		Target:  target,
		Payload: payload,
	}

	// build the payload
	if res, err = s.buildPayload(&build); err != nil {
		return "", err
	}

	// build the stage
	if err = s.buildStage(&build); err != nil {
		return "", err
	}

	s.Builds = append(s.Builds, build.ID)
	return res, nil
}
