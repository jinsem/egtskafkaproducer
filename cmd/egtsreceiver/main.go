package main

import (
	producer2 "github.com/jinsem/egtskafkaproducer/pkg/common"
	"github.com/labstack/gommon/log"
	"net"
	"os"
)

var (
	settings producer2.Settings
	logger   *log.Logger
)

const maxArgCnt = 2

func main() {
	logger = log.New("-")
	logger.SetHeader("${time_rfc3339_nano} ${short_file}:${line} ${level} -${message}")
	loadSettings()
	logger.SetLevel(settings.Log.GetLevel())
	producer := producer2.KafkaProducer{}
	if e := producer.Initialize(&settings.Kafka, logger); e != nil {
		logger.Fatalf("Persister initialization failed: %v", e)
	}
	defer func() {
		err := producer.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	startTcpListener(settings.App.GetFullAddress(), producer)

}

func loadSettings() {
	if len(os.Args) == maxArgCnt {
		if err := settings.LoadFromFile(os.Args[1]); err != nil {
			logger.Fatalf("Ошибка загрузки файла настроек: %v", err)
		}
	} else {
		if err := settings.LoadFromEnv(); err != nil {
			logger.Fatalf("Ошибка настройки приложения: %v", err)
		}
	}
}

func startTcpListener(srvAddress string, producer producer2.KafkaProducer) {
	listener, err := net.Listen("tcp", srvAddress)
	if err != nil {
		logger.Fatalf("Cannot open TCP connection: %v", err)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Infof("Listener is running on %s...", srvAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("Connection error: %v", err)
		} else {
			go handleReceivedPackage(conn, producer)
		}
	}
}
