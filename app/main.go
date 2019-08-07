package main

import (
	"github.com/labstack/gommon/log"
	"os"
)

var (
	settings Settings
	logger   *log.Logger
)

const MAX_ARG_CNT = 2

func main() {
	loadSettings()
	initLogger()
}

func loadSettings() {
	if len(os.Args) == MAX_ARG_CNT {
		if err := settings.Load(os.Args[1]); err != nil {
			logger.Fatalf("Application configuration cannot be parsed: %v", err)
		}
	} else {
		logger.Fatalf("Path to configuration is not set")
	}
}

func initLogger() {
	logger = log.New("-")
	logger.SetHeader("${time_rfc3339_nano} ${short_file}:${line} ${level} -${message}")
	logger.SetLevel(settings.Log.getLevel())
}
