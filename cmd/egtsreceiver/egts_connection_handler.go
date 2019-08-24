// This is modified version of this file https://github.com/kuznetsovin/egts/blob/master/cmd/receiver/egts_handler.go
package main

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
	"github.com/jinsem/egtskafkaproducer/pkg/common"
	egts "github.com/kuznetsovin/egts/pkg/egtslib"
	"io"
	"net"
	"time"
)

const (
	egtsPcOk        = 0
	headerLen       = 10
	srcLen          = 2
	protocolVersion = 0x01
	existsFlag      = "1"
)

func handleReceivedPackage(conn net.Conn, producer common.KafkaProducer) {

	var (
		readyToPersist    bool
		srResultCodePkg   []byte
		serviceType       uint8
		srResponsesRecord egts.RecordDataSet
		recvPacket        []byte
		deviceImei        string
	)
	logger.Debugf("Соединение установлено. Адрес устройства: %s", conn.RemoteAddr())
	for {
	Received:
		serviceType = 0
		srResponsesRecord = nil
		srResultCodePkg = nil
		recvPacket = nil

		connTimer := time.NewTimer(settings.App.GetConnectionTimeToLiveSec())

		headerBuf := make([]byte, headerLen)

		_, err := conn.Read(headerBuf)

		switch err {
		case nil:
			connTimer.Reset(settings.App.GetConnectionTimeToLiveSec())

			if headerBuf[0] != protocolVersion {
				_ = conn.Close()
				logger.Warnf("Полученный пакет не соответствует спецификации ЕГТС. Соединение с %s будет закрыто", conn.RemoteAddr())
				return
			}

			// packageLen = header length (HL) + body length (FDL) + package CRC (2 bytes) if FDL is not empty (see order №285)
			bodyLen := binary.LittleEndian.Uint16(headerBuf[5:7])
			pkgLen := uint16(headerBuf[3])
			if bodyLen > 0 {
				pkgLen += bodyLen + srcLen
			}
			buf := make([]byte, pkgLen-headerLen)
			if _, err := io.ReadFull(conn, buf); err != nil {
				logger.Errorf("Ошибка чтения тела пакета: %v", err)
				_ = conn.Close()
				return
			}
			recvPacket = append(headerBuf, buf...)
		case io.EOF:
			<-connTimer.C
			_ = conn.Close()
			logger.Warnf("Соединение с %s закрыто по таймауту", conn.RemoteAddr())
			return
		default:
			logger.Errorf("Ошибка при получении: %v", err)
			_ = conn.Close()
			return
		}

		logger.Debugf("Принят пакет: %X\v", recvPacket)
		pkg := egts.Package{}
		resultCode, err := pkg.Decode(recvPacket)
		if err != nil {
			logger.Warn("Ошибка расшифровки пакета")
			logger.Error(err)

			resp, err := createPtResponse(&pkg, resultCode, serviceType, nil)
			if err != nil {
				logger.Errorf("Ошибка сборки ответа EGTS_PT_RESPONSE: %v", err)
				goto Received
			}
			_, _ = conn.Write(resp)

			goto Received
		}

		switch pkg.PacketType {
		case egts.PtAppdataPacket:
			logger.Info("Тип пакета EGTS_PT_APPDATA")

			for _, rec := range *pkg.ServicesFrameData.(*egts.ServiceDataSet) {
				exportPacket := egtsschema.EgtsPackage{
					AnalogSensors:    &egtsschema.UnionArrayAnalogSensorNull{},
					LiquidSensors:    &egtsschema.UnionArrayLiquidSensorNull{},
					Latitude:         &egtsschema.UnionNullDouble{},
					Longitude:        &egtsschema.UnionNullDouble{},
					Speed:            &egtsschema.UnionNullInt{},
					Direction:        &egtsschema.UnionNullInt{},
					NumOfSatelites:   &egtsschema.UnionNullInt{},
					Pdop:             &egtsschema.UnionNullInt{},
					Hdop:             &egtsschema.UnionNullInt{},
					Vdop:             &egtsschema.UnionNullInt{},
					NavigationSystem: egtsschema.NavigationSystem(egtsschema.NavigationSystemUknown),
				}
				readyToPersist = false
				packetIDBytes := make([]byte, 4)

				srResponsesRecord = append(srResponsesRecord, egts.RecordData{
					SubrecordType:   egts.SrRecordResponseType,
					SubrecordLength: 3,
					SubrecordData: &egts.SrResponse{
						ConfirmedRecordNumber: rec.RecordNumber,
						RecordStatus:          egtsPcOk,
					},
				})
				serviceType = rec.SourceServiceType
				logger.Info("Тип сервиса ", serviceType)

				exportPacket.ClinetId = int64(rec.ObjectIdentifier)
				for _, subRec := range rec.RecordDataSet {
					switch subRecData := subRec.SubrecordData.(type) {
					case *egts.SrTermIdentity:
						logger.Debugf("Разбор подзаписи EGTS_SR_TERM_IDENTITY")
						deviceImei = subRecData.IMEI
						if srResultCodePkg, err = createSrResultCode(&pkg, egtsPcOk); err != nil {
							logger.Errorf("Ошибка сборки EGTS_SR_RESULT_CODE: %v", err)
						}
					case *egts.SrAuthInfo:
						logger.Debugf("Разбор подзаписи EGTS_SR_AUTH_INFO")
						if srResultCodePkg, err = createSrResultCode(&pkg, egtsPcOk); err != nil {
							logger.Errorf("Ошибка сборки EGTS_SR_RESULT_CODE: %v", err)
						}
					case *egts.SrResponse:
						logger.Debugf("Разбор подзаписи EGTS_SR_RESPONSE")
						goto Received
					case *egts.SrPosData:
						logger.Debugf("Разбор подзаписи EGTS_SR_POS_DATA")
						readyToPersist = true
						setSrPosData(&exportPacket, subRecData)
					case *egts.SrExtPosData:
						logger.Debugf("Разбор подзаписи EGTS_SR_EXT_POS_DATA")
						setSrExtPosData(&exportPacket, subRecData)
					case *egts.SrAdSensorsData:
						logger.Debugf("Разбор подзаписи EGTS_SR_AD_SENSORS_DATA")
						readyToPersist = true
						setSrAdSensorsData(&exportPacket, subRecData)
					case *egts.SrAbsCntrData:
						logger.Debugf("Разбор подзаписи EGTS_SR_ABS_CNTR_DATA")
						readyToPersist = true
						switch subRecData.CounterNumber {
						case 110:
							// Три младших байта номера передаваемой записи (идет вместе с каждой POS_DATA).
							binary.BigEndian.PutUint32(packetIDBytes, subRecData.CounterValue)
							exportPacket.PacketID = int64(subRecData.CounterValue)
						case 111:
							// один старший байт номера передаваемой записи (идет вместе с каждой POS_DATA).
							tmpBuf := make([]byte, 4)
							binary.BigEndian.PutUint32(tmpBuf, subRecData.CounterValue)
							if len(packetIDBytes) == 4 {
								packetIDBytes[3] = tmpBuf[3]
							} else {
								packetIDBytes = tmpBuf
							}
							exportPacket.PacketID = int64(binary.LittleEndian.Uint32(packetIDBytes))
						}
					case *egts.SrLiquidLevelSensor:
						logger.Debugf("Разбор подзаписи EGTS_SR_LIQUID_LEVEL_SENSOR")
						readyToPersist = true
						setSrLiquidLevelSensor(&exportPacket, subRecData)
					}
				}

				if readyToPersist {
					exportPacket.Guid = fmt.Sprintf("%s", uuid.New())
					exportPacket.Imei = &egtsschema.UnionNullString{String: deviceImei}
					if err := producer.Produce(&exportPacket); err != nil {
						logger.Error(err)
					}
				}
			}

			resp, err := createPtResponse(&pkg, resultCode, serviceType, srResponsesRecord)
			if err != nil {
				logger.Errorf("Ошибка сборки ответа: %v", err)
				goto Received
			}

			_, err = conn.Write(resp)
			if err == nil {
				logger.Debugf("Отправлен пакет EGTS_PT_RESPONSE: %X", resp)
			} else {
				logger.Errorf("Ошибка отправки пакета EGTS_PT_RESPONSE: %v", err)
				goto Received
			}

			if len(srResultCodePkg) > 0 {
				_, err = conn.Write(srResultCodePkg)
				if err == nil {
					logger.Debugf("Отправлен пакет EGTS_SR_RESULT_CODE: %X", resp)
				} else {
					logger.Errorf("Ошибка отправки пакета EGTS_SR_RESULT_CODE: %v", err)
				}
			}
		case egts.PtResponsePacket:
			logger.Debug("Тип пакета EGTS_PT_RESPONSE")
		}
	}
}

func setSrPosData(exportPacket *egtsschema.EgtsPackage, subRecData *egts.SrPosData) {
	exportPacket.MeasurementTimestamp = subRecData.NavigationTime.Unix()
	exportPacket.ReceivedTimestamp = time.Now().UTC().Unix()
	setUnionNullDoubleVal(subRecData.Latitude, exportPacket.Latitude)
	setUnionNullDoubleVal(subRecData.Longitude, exportPacket.Longitude)
	setUnionNullInt(int32(subRecData.Speed), exportPacket.Speed)
	setUnionNullInt(int32(subRecData.Direction), exportPacket.Direction)
}

func setSrExtPosData(exportPacket *egtsschema.EgtsPackage, subRecData *egts.SrExtPosData) {
	setUnionNullInt(int32(subRecData.Satellites), exportPacket.NumOfSatelites)
	setUnionNullInt(int32(subRecData.PositionDilutionOfPrecision), exportPacket.Pdop)
	setUnionNullInt(int32(subRecData.HorizontalDilutionOfPrecision), exportPacket.Hdop)
	setUnionNullInt(int32(subRecData.VerticalDilutionOfPrecision), exportPacket.Vdop)
	exportPacket.NavigationSystem = toNavigationSystem(subRecData.NavigationSystem)
}

func setSrAdSensorsData(exportPacket *egtsschema.EgtsPackage, subRecData *egts.SrAdSensorsData) {
	if subRecData.AnalogSensorFieldExists1 == existsFlag {
		sensor := egtsschema.AnalogSensor{1, int32(subRecData.AnalogSensor1)}
		exportPacket.AnalogSensors.ArrayAnalogSensor = append(exportPacket.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
	if subRecData.AnalogSensorFieldExists2 == existsFlag {
		sensor := egtsschema.AnalogSensor{2, int32(subRecData.AnalogSensor2)}
		exportPacket.AnalogSensors.ArrayAnalogSensor = append(exportPacket.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
	if subRecData.AnalogSensorFieldExists3 == existsFlag {
		sensor := egtsschema.AnalogSensor{3, int32(subRecData.AnalogSensor3)}
		exportPacket.AnalogSensors.ArrayAnalogSensor = append(exportPacket.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
	if subRecData.AnalogSensorFieldExists4 == existsFlag {
		sensor := egtsschema.AnalogSensor{4, int32(subRecData.AnalogSensor4)}
		exportPacket.AnalogSensors.ArrayAnalogSensor = append(exportPacket.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
	if subRecData.AnalogSensorFieldExists5 == existsFlag {
		sensor := egtsschema.AnalogSensor{5, int32(subRecData.AnalogSensor5)}
		exportPacket.AnalogSensors.ArrayAnalogSensor = append(exportPacket.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
	if subRecData.AnalogSensorFieldExists6 == existsFlag {
		sensor := egtsschema.AnalogSensor{6, int32(subRecData.AnalogSensor6)}
		exportPacket.AnalogSensors.ArrayAnalogSensor = append(exportPacket.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
	if subRecData.AnalogSensorFieldExists7 == existsFlag {
		sensor := egtsschema.AnalogSensor{7, int32(subRecData.AnalogSensor7)}
		exportPacket.AnalogSensors.ArrayAnalogSensor = append(exportPacket.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
	if subRecData.AnalogSensorFieldExists8 == existsFlag {
		sensor := egtsschema.AnalogSensor{8, int32(subRecData.AnalogSensor8)}
		exportPacket.AnalogSensors.ArrayAnalogSensor = append(exportPacket.AnalogSensors.ArrayAnalogSensor, &sensor)
	}
}

func setSrLiquidLevelSensor(exportPacket *egtsschema.EgtsPackage, subRecData *egts.SrLiquidLevelSensor) {
	valueMillimetres := int32(0)
	valueLitres := int32(0)
	switch subRecData.LiquidLevelSensorValueUnit {
	case "00", "01":
		valueMillimetres = int32(subRecData.LiquidLevelSensorData)
	case "10":
		valueLitres = int32(subRecData.LiquidLevelSensorData * 10)
	}
	sensor := egtsschema.LiquidSensor{
		int32(subRecData.LiquidLevelSensorNumber),
		subRecData.LiquidLevelSensorErrorFlag,
		valueMillimetres,
		valueLitres}
	exportPacket.LiquidSensors.ArrayLiquidSensor = append(exportPacket.LiquidSensors.ArrayLiquidSensor, &sensor)
}
