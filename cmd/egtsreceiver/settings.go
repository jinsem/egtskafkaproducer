package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

const prefix = "EGTS_RECEIVER"
const appHostname = "APP_HOSTNAME"
const appPort = "APP_PORT"
const appConnectiontimetolivesec = "APP_CONNECTIONTIMETOLIVESEC"
const logLevel = "LOG_LEVEL"
const kafkaBrokers = "KAFKA_BROKERS"
const kafkaOutputtopicname = "KAFKA_OUTPUTTOPICNAME"

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

func (s *Settings) Load(configPath string) error {

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		_, confNotFound := err.(viper.ConfigFileNotFoundError)
		_, pathErr := err.(*os.PathError)
		if !confNotFound && !pathErr {
			return fmt.Errorf("Configuration file cannot be loaded because of the error: %v \n", err)
		}
	}

	if err := viper.Unmarshal(s); err != nil {
		return fmt.Errorf("Ошибка создания объекта настроек %v", err)
	}

	if err := updateFromEnv(s); err != nil {
		return fmt.Errorf("Ошибка инициализации объекта настроек %v", err)
	}
	return nil
}

func updateFromEnv(s *Settings) error {
	envS := Settings{
		App: AppSettings{
			HostName:                viper.GetString(appHostname),
			Port:                    viper.GetString(appPort),
			ConnectionTimeToLiveSec: viper.GetInt(appConnectiontimetolivesec),
		},
		Log: LogSettings{
			Level: viper.GetString(logLevel),
		},
		Kafka: KafkaSettings{
			Brokers:         strings.Split(viper.GetString(kafkaBrokers), ","),
			OutputTopicName: viper.GetString(kafkaOutputtopicname),
		},
	}
	return mergo.Merge(s, &envS)
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
