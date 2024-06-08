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

type StageConfig struct {
  Target    *Target 
  ID        string
  Address   string
  Port      int
}

var (
  STAGE_CONFIG = path.Join("src", "config.h")
  STAGE_OUTPUT = "stage" 
)

var stages []string

func StageLoad() error {
  entries, err := os.ReadDir(global.CONFIG_DATADIR)
  if os.IsNotExist(err) {
    return nil
  }

  if err != nil {
    return err
  }

  for _, e := range entries {
    if strings.HasSuffix(e.Name(), "_stage") {
      stages = append(stages, e.Name())
      break
    }
  }

  return nil
}

func StageGet(id string) string {
  if len(id) != global.ID_LEN {
    return ""
  }

  for _, s := range stages {
    fp := path.Join(global.CONFIG_DATADIR, s)
    if strings.HasPrefix(s, id) && util.Access(fp){
      return fp 
    }
  }

  return ""
}

func StageHas(id string) bool {
  return StageGet(id) != ""
}

func StageAdd(id string, src string) (string, error) {
  name := fmt.Sprintf("%s_stage", id)
  
  err := os.Mkdir(global.CONFIG_DATADIR, 0777)
  if err != nil && !os.IsExist(err) {
    return "", err
  }

  err = util.CopyFile(path.Join(global.CONFIG_DATADIR, name), src)
  if err != nil {
    return "", err
  }

  stages = append(stages, name)
  return name, nil
}

func StageBuild(conf StageConfig) (string, error) {
  stagedir := path.Join(global.CONFIG_PAYLOADDIR, "stage")
  builddir := path.Join("/tmp", "ezcat_stage_"+conf.ID)

  err := util.CopySimple(builddir, stagedir)
  if err != nil {
    return "", fmt.Errorf("failed to copy stage dir: %s", err.Error()) 
  }

  cfgfile := path.Join(builddir, STAGE_CONFIG)
  err = os.WriteFile(cfgfile, []byte(fmt.Sprintf(`#pragma once
#include <stdbool.h>

#define SERVER_ADDRESS "%s"
#define SERVER_PORT %d
#define ID "%s"
#define DEBUG %s`, 
  conf.Address, 
  conf.Port, 
  conf.ID, 
  util.BoolToStr(global.CONFIG_DEBUG))), os.ModePerm)

  if err != nil {
    return "", fmt.Errorf("failed to write to the config file: %s", err.Error()) 
  }

  out, err := util.RunBuild(builddir, conf.Target.CodeName())

  if err != nil {
    log.Err("Failed to run the build command for the stage:\n=====\n%s\n=====", out)
    return "", fmt.Errorf("failed to run the build: %s", err.Error()) 
  }

  path, err := StageAdd(conf.ID, path.Join(builddir, STAGE_OUTPUT))
  if err != nil {
    return "", fmt.Errorf("failed to add stage build: %s", err.Error())
  }

  return path, nil 
}
