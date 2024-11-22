package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/ngn13/ezcat/server/util"
)

type Struct struct {
	Version    string // not modifiable
	HTTP_Port  int    `env:"HTTP_PORT"`
	C2_Port    int    `env:"C2_PORT"`
	Password   string `env:"PASSWORD"`
	StaticDir  string `env:"STATIC_DIR"`
	DistDir    string `env:"DIST_DIR"`
	PayloadDir string `env:"PAYLOAD_DIR"`
	ShellIP    string `env:"SHELLIP"`
	Megamind   bool   `env:"MEGAMIND"`
	Debug      bool   `env:"DEBUG"`
}

func (conf *Struct) Load() error {
	var (
		val reflect.Value
		typ reflect.Type

		field_val reflect.Value
		field_typ reflect.StructField

		env_name string
		env_val  string

		err error
		ok  bool
	)

	val = reflect.ValueOf(conf).Elem()
	typ = val.Type()

	for i := 0; i < val.NumField(); i++ {
		field_val = val.Field(i)
		field_typ = typ.Field(i)

		if env_name, ok = field_typ.Tag.Lookup("env"); !ok || !field_val.CanSet() {
			continue
		}

		env_name = fmt.Sprintf("EZCAT_%s", env_name)

		if env_val = os.Getenv(env_name); env_val == "" {
			continue
		}

		switch field_val.Kind() {
		case reflect.String:
			field_val.SetString(env_val)

		case reflect.Int:
			var env_val_int int
			if env_val_int, err = strconv.Atoi(env_val); err != nil {
				return fmt.Errorf("%s should be an integer", env_name)
			}
			field_val.SetInt(int64(env_val_int))

		case reflect.Bool:
			var env_val_bool bool
			if env_val == "true" || env_val == "1" {
				env_val_bool = true
			} else if env_val == "false" || env_val == "0" {
				env_val_bool = false
			} else {
				return fmt.Errorf("%s should be a boolean", env_name)
			}
			field_val.SetBool(env_val_bool)
		}
	}

	return nil
}

func New() (*Struct, error) {
	var (
		conf Struct
		err  error
	)

	conf.Version = "2.5"
	conf.HTTP_Port = 5566
	conf.C2_Port = 5567
	conf.Password = "ezcat"
	conf.DistDir = "./data"
	conf.StaticDir = ""
	conf.PayloadDir = "../payloads"
	conf.Megamind = true
	conf.Debug = false

	if err = conf.Load(); err != nil {
		return nil, err
	}

	if !util.IsValidPort(conf.HTTP_Port) {
		return nil, fmt.Errorf("invalid port number for the HTTP port: %d", conf.HTTP_Port)
	}

	if !util.IsValidPort(conf.C2_Port) {
		return nil, fmt.Errorf("invalid port number for the C2 port: %d", conf.C2_Port)
	}

	return &conf, nil
}
