package config

import (
	"log"

	"github.com/spf13/viper"
)

func Load(configName string) (config Configuration) {
	viper.SetConfigName(configName)
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return config
}
