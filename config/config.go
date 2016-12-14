package config

import (
	"flag"
	"fmt"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

var GlobalConfig Config

type Config struct {
	ServAddr     string `toml:"serv_addr"`
	DSN          string `toml:"dsn"`
	FrontendRoot string `toml:"frontend_root"`
}

func init() {
	var fpath string
	flag.StringVar(&fpath, "c", "config.toml", "Configuration file to use")
	flag.Parse()
	fmt.Println(fpath)
	_, err := toml.DecodeFile(fpath, &GlobalConfig)
	if err != nil {
		panic(err)
	}
	log.Infof("config module inited")
}
