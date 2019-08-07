package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type appSettings struct {
	HostName                string
	Port                    string
	ConnectionTimeToLiveSec int
}

func (c *appSettings) Load(configPath string) error {

	viper.SetConfigName(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("Configuration file %s not found \n", configPath)
		} else {
			return fmt.Errorf("Configuration file cannot be loaded because of the error: %v \n", err)
		}
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		return fmt.Errorf("Unable to decode configuration %v", err)
	}
	return nil
}
