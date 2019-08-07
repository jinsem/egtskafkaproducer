package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
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
	appS := AppSettings{}
	logS := LogSettings{}
	kafkaS := KafkaSettings{}
	err := viper.UnmarshalKey(AppKey, &appS)
	if err != nil {
		return fmt.Errorf("Unable to decode application configuration %v", err)
	}
	err = viper.UnmarshalKey(LogKey, &logS)
	if err != nil {
		return fmt.Errorf("Unable to decode log configuration %v", err)
	}
	err = viper.UnmarshalKey(KafkaKey, &kafkaS)
	if err != nil {
		return fmt.Errorf("Unable to decode Kafka configuration %v", err)
	}
	c.App = appS
	c.Log = logS
	c.Kafka = kafkaS
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
