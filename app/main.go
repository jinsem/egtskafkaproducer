package main

import (
	"github.com/labstack/gommon/log"
	"net"
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
	producer := EgtsProducer{}
	if e := producer.Initialize(settings.Kafka); e != nil {
		logger.Fatalf("Producer initialization failed: %v", e)
	}
	defer producer.Close()
	startTcpListener(settings.App.getFullAddress(), producer)

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

func startTcpListener(srvAddress string, producer EgtsProducer) {
	listener, err := net.Listen("tcp", srvAddress)
	if err != nil {
		logger.Fatalf("Cannot open TCP connection: %v", err)
	}
	defer listener.Close()

	logger.Infof("Listener is running on %s...", srvAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("Connection error: %v", err)
		} else {
			go handleRecvPkg(conn, producer)
		}
	}
}
