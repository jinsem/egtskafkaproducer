package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const prefix = "EGTS_RECEIVER"
const appHostname = "APP_HOSTNAME"
const appPort = "APP_PORT"
const appConnectiontimetolivesec = "APP_CONNECTIONTIMETOLIVESEC"
const logLevel = "LOG_LEVEL"
const kafkaBrokers = "KAFKA_BROKERS"
const kafkaSchemaRegistryUrl = "KAFKA_SCHEMAREGISTRYURL"
const kafkaOutputtopicname = "KAFKA_OUTPUTTOPICNAME"
const defaultTtl = 60

type Settings struct {
	App   AppSettings   `yaml:"app"`
	Log   LogSettings   `yaml:"log"`
	Kafka KafkaSettings `yaml:"kafka"`
}

type AppSettings struct {
	HostName                string `yaml:"HostName"`
	Port                    string `yaml:"Port"`
	ConnectionTimeToLiveSec int    `yaml:"ConnectionTimeToLiveSec"`
}

type LogSettings struct {
	Level string `yaml:"Level"`
}

type KafkaSettings struct {
	Brokers           []string `yaml:"Brokers"`
	SchemaRegistryUrl string   `yaml:"SchemaRegistryUrl"`
	OutputTopicName   string   `yaml:"OutputTopicName"`
}

func (s *Settings) LoadFromFile(configPath string) error {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err == nil {
		if err = yaml.Unmarshal(yamlFile, s); err != nil {
			return fmt.Errorf("Ошибка создания объекта настроек %v", err)
		}
	} else {
		log.Warnf("Файл настроек не найден или не может быть загружен: #%v ", err)
	}
	if err := validateSettings(s); err != nil {
		return fmt.Errorf("Ошибка валидации настроек программы %v", err)
	}
	return nil
}

func validateSettings(s *Settings) error {
	if strings.Trim(s.App.Port, " ") == "" {
		return fmt.Errorf("Не задано значение порта приложения")
	}
	if s.App.ConnectionTimeToLiveSec <= 0 {
		return fmt.Errorf("Время поддержания соединения должно быть положительным числом ")
	}
	if len(s.Kafka.Brokers) == 0 {
		return fmt.Errorf("Не заданы адреса брокеров Kafka")
	}
	if strings.Trim(s.Kafka.OutputTopicName, " ") == "" {
		return fmt.Errorf("Не задано имя топика для публикации данных")
	}
	return nil
}

func (s *Settings) LoadFromEnv() error {
	if err := updateFromEnv(s); err != nil {
		return fmt.Errorf("Ошибка создания объекта настроек %v", err)
	}
	if err := validateSettings(s); err != nil {
		return fmt.Errorf("Ошибка валидации настроек программы %v", err)
	}
	return nil
}

func updateFromEnv(s *Settings) error {
	ttl, err := strconv.Atoi(os.Getenv(makeKey(appConnectiontimetolivesec)))
	if err != nil {
		ttl = defaultTtl
	}
	envS := Settings{
		App: AppSettings{
			HostName:                os.Getenv(makeKey(appHostname)),
			Port:                    os.Getenv(makeKey(appPort)),
			ConnectionTimeToLiveSec: ttl,
		},
		Log: LogSettings{
			Level: os.Getenv(makeKey(logLevel)),
		},
		Kafka: KafkaSettings{
			Brokers:           strings.Split(os.Getenv(makeKey(kafkaBrokers)), ","),
			SchemaRegistryUrl: os.Getenv(makeKey(kafkaSchemaRegistryUrl)),
			OutputTopicName:   os.Getenv(makeKey(kafkaOutputtopicname)),
		},
	}
	return mergo.Merge(s, &envS)
}

func makeKey(suffix string) string {
	return prefix + "_" + suffix
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
