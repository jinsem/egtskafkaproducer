package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"time"
)

const AppKey = "app"
const LogKey = "log"
const KafkaKey = "kafka"

type Settings struct {
	App   AppSettings
	Log   LogSettings
	Kafka KafkaSettings
}

type AppSettings struct {
	HostName                string
	Port                    string
	ConnectionTimeToLiveSec int
}

type LogSettings struct {
	Level string
}

type KafkaSettings struct {
	Brokers         []string
	OutputTopicName string
}

func (c *Settings) Load(configPath string) error {

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("Configuration file %s not found \n", configPath)
		} else {
			return fmt.Errorf("Configuration file cannot be loaded because of the error: %v \n", err)
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("Unable to decode configuration %v", err)
	}
	return nil
}

func (l *LogSettings) getLevel() log.Lvl {
	var lvl log.Lvl

	switch l.Level {
	case "DEBUG":
		lvl = log.DEBUG
		break
	case "INFO":
		lvl = log.INFO
		break
	case "WARN":
		lvl = log.WARN
		break
	case "ERROR":
		lvl = log.ERROR
		break
	default:
		lvl = log.INFO
	}
	return lvl
}

func (a *AppSettings) getFullAddress() string {
	return a.HostName + ":" + a.Port
}

func (s *AppSettings) getConnectionTimeToLiveSec() time.Duration {
	return time.Duration(s.ConnectionTimeToLiveSec) * time.Second
}
