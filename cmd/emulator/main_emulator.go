package main

import (
	"fmt"
	"github.com/google/uuid"
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
	"github.com/jinsem/egtskafkaproducer/pkg/common"
	"github.com/labstack/gommon/log"
	"time"
)

var (
	logger              *log.Logger
	allowedClientIds    = []int64{1, 2, 3, 4}
	allowedInputNumbers = []int64{1, 2, 3, 4}
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
	for i := 0; i < 2; i++ {
		fakePackege := makEgtsPackage(i)
		_ = kafkaProducer.Produce(&fakePackege)
	}
	_ = kafkaProducer.Close()
}

func makEgtsPackage(i int) egtsschema.EgtsPackage {
	result := egtsschema.EgtsPackage{
		AnalogSensors: &egtsschema.UnionArrayAnalogSensorNull{},
		LiquidSensors: &egtsschema.UnionArrayLiquidSensorNull{},
	}
	result.ClinetId = int64(i)
	result.PacketID = int64(i)
	result.Guid = fmt.Sprintf("%s", uuid.New())
	result.Imei = fmt.Sprintf("imei %d", i)

	result.MeasurementTimestamp = 12345
	result.ReceivedTimestamp = time.Now().UTC().Unix()
	result.Latitude = 12.9
	result.Longitude = 11.5
	result.Speed = int32(50)
	result.Direction = int32(55)
	result.NumOfSatelites = int32(30)
	result.Pdop = int32(11)
	result.Hdop = int32(12)
	result.Vdop = int32(13)
	result.NavigationSystem = egtsschema.NavigationSystem(egtsschema.NavigationSystemGLONASS)

	sensor1 := egtsschema.AnalogSensor{1, int32(100)}
	result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor1)
	sensor2 := egtsschema.AnalogSensor{2, int32(200)}
	result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor2)
	sensor3 := egtsschema.AnalogSensor{3, int32(300)}
	result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor3)
	sensor4 := egtsschema.AnalogSensor{4, int32(400)}
	result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor4)
	sensor5 := egtsschema.AnalogSensor{5, int32(500)}
	result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor5)
	sensor6 := egtsschema.AnalogSensor{6, int32(600)}
	result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor6)
	sensor7 := egtsschema.AnalogSensor{7, int32(700)}
	result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor7)
	sensor8 := egtsschema.AnalogSensor{8, int32(800)}
	result.AnalogSensors.ArrayAnalogSensor = append(result.AnalogSensors.ArrayAnalogSensor, &sensor8)
	return result
}
