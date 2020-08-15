package config

import (
	"github.com/tkanos/gonfig"

	"../utils"
)

type configuration struct {
	Prefix    string
	AuthToken string
	Hostname  string
	Port      int
}

// Configuration has all the configuration variables stored in it
var Configuration = configuration{}

// InitConfig initializes the Configuration object
func InitConfig() {
	err := gonfig.GetConf("./config.json", &Configuration)
	utils.Check(err)
}
