package config

import (
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	log.Print("This is the environment: ", env)
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("config/")
	config.AddConfigPath("./internal/api/config/")
	config.AddConfigPath(".")             // used for docker
	config.AddConfigPath("../../config/") // used for unit tests
	err = config.ReadInConfig()
	if err != nil {
		log.Fatalf("error on parsing configuration file: %s", err)
	}
}

func GetConfig() *viper.Viper {
	return config
}
