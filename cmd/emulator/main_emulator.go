package main

import (
	"fmt"
	"github.com/google/uuid"
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
	"github.com/jinsem/egtskafkaproducer/pkg/common"
	"github.com/labstack/gommon/log"
	"math/rand"
	"time"
)

var (
	logger             *log.Logger
	allowedClientIds   = []int64{1, 2, 3, 4}
	allowedClientImeis = []string{"111111111", "222222222222", "333333333333", "444444444444"}
	maxInputNumber     = 8
)

func main() {
	logger = log.New("-")
	logger.SetHeader("${time_rfc3339_nano} ${short_file}:${line} ${level} -${message}")
	logger.SetLevel(log.DEBUG)
	settings := common.KafkaSettings{
		Brokers:           []string{"localhost:29092"},
		SchemaRegistryUrl: "http://localhost:8081",
		OutputTopicName:   "sensor_data",
	}
	kafkaProducer := common.KafkaProducer{}
	_ = kafkaProducer.Initialize(&settings, logger)
	for i := 0; i < 1; i++ {
		fakePackege := makMeasurementPackage()
		_ = kafkaProducer.Produce(&fakePackege)
	}
	_ = kafkaProducer.Close()
}

func makMeasurementPackage() egtsschema.MeasurementPackage {
	result := egtsschema.MeasurementPackage{
		AnalogSensors: &egtsschema.UnionArrayAnalogSensorNull{},
		LiquidSensors: &egtsschema.UnionArrayLiquidSensorNull{},
	}
	rndClientId, rndImei := getRandomClientIdAndImei()
	result.ClinetId = rndClientId
	result.PacketID = time.Now().UnixNano()
	result.Guid = fmt.Sprintf("%s", uuid.New())
	result.Imei = &egtsschema.UnionNullString{String: rndImei}

	curTimeStamp := time.Now().UTC().Unix()
	result.MeasurementTimestamp = curTimeStamp
	result.ReceivedTimestamp = curTimeStamp - 1000
	result.Latitude = &egtsschema.UnionNullDouble{Double: 12.9, UnionType: egtsschema.UnionNullDoubleTypeEnumDouble}
	result.Longitude = &egtsschema.UnionNullDouble{Double: 11.5, UnionType: egtsschema.UnionNullDoubleTypeEnumDouble}
	result.Speed = &egtsschema.UnionNullInt{Int: int32(0), UnionType: egtsschema.UnionNullIntTypeEnumInt}
	result.Direction = &egtsschema.UnionNullInt{Int: int32(0), UnionType: egtsschema.UnionNullIntTypeEnumInt}
	result.NumOfSatelites = &egtsschema.UnionNullInt{Int: int32(10), UnionType: egtsschema.UnionNullIntTypeEnumInt}
	result.Pdop = &egtsschema.UnionNullInt{Int: int32(1), UnionType: egtsschema.UnionNullIntTypeEnumInt}
	result.Hdop = &egtsschema.UnionNullInt{Int: int32(2), UnionType: egtsschema.UnionNullIntTypeEnumInt}
	result.Vdop = &egtsschema.UnionNullInt{Int: int32(3), UnionType: egtsschema.UnionNullIntTypeEnumInt}
	result.NavigationSystem = egtsschema.NavigationSystem(egtsschema.NavigationSystemGLONASS)

	for i := 1; i <= maxInputNumber; i++ {
		sensor := egtsschema.AnalogSensor{int32(i), getRandomInpitValue()}
		result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
	return result
}

func getRandomClientIdAndImei() (int64, string) {
	rndIdx := rand.Intn(len(allowedClientIds))
	return allowedClientIds[rndIdx], allowedClientImeis[rndIdx]
}

func getRandomInpitValue() int32 {
	return rand.Int31n(10000)
}
