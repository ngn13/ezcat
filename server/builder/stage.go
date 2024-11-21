package builder

import (
	"fmt"
	"os"
	"path"

	"github.com/ngn13/ezcat/server/log"
	"github.com/ngn13/ezcat/server/util"
)

var (
	STAGE_CONFIG = path.Join("src", "config.h")
	STAGE_OUTPUT = "stage"
)

func (s *Struct) buildStage(b *build) (err error) {
	stagedir := path.Join(s.Config.PayloadDir, "stage")
	builddir := path.Join("/tmp", "ezcat_stage_"+b.ID)

	if err = util.CopySimple(builddir, stagedir); err != nil {
		return fmt.Errorf("failed to copy stage dir: %s", err.Error())
	}

	options := map[string]string{
		"STAGE_ID":          b.ID,
		"STAGE_SERVER_HOST": b.Host,
		"STAGE_SERVER_PORT": fmt.Sprintf("%d", s.Config.C2_Port),
	}

	if s.Config.Debug {
		options["STAGE_DEBUG"] = "1"
	} else {
		options["STAGE_DEBUG"] = "0"
	}

	if out, err := util.RunBuild(builddir, b.Target.CodeName(), options); err != nil {
		log.Fail("failed to run the build command for the stage:\n=====\n%s\n=====", out)
		return fmt.Errorf("failed to run the build: %s", err.Error())
	}

	if err = os.Mkdir(s.Config.DistDir, 0777); err != nil && !os.IsExist(err) {
		return err
	}

	dst := path.Join(s.Config.DistDir, b.ID)
	src := path.Join(builddir, STAGE_OUTPUT)

	if err = util.CopyFile(dst, src); err != nil {
		return err
	}

	return nil
}
