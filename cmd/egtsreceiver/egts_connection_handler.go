// This is modified version of this file https://github.com/kuznetsovin/egts/blob/master/cmd/receiver/egts_handler.go
package main

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
	egts "github.com/kuznetsovin/egts/pkg/egtslib"
	"io"
	"net"
	"time"
)

const (
	egtsPcOk        = 0
	headerLen       = 10
	protocolVersion = 0x01
)

func handleReceivedvPackage(conn net.Conn, producer EgtsKafkaPersister) {

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

		connTimer := time.NewTimer(settings.App.getConnectionTimeToLiveSec())

		headerBuf := make([]byte, headerLen)

		_, err := conn.Read(headerBuf)

		switch err {
		case nil:
			connTimer.Reset(settings.App.getConnectionTimeToLiveSec())

			if headerBuf[0] != protocolVersion {
				_ = conn.Close()
				logger.Warnf("Полученный пакет не соответствует спецификации ЕГТС. Соединение с %s будет закрыто", conn.RemoteAddr())
				return
			}

			// packageLen = header length (HL) + body length (FDL) + package CRC (2 bytes) if FDL is not empty (see order №285)
			bodyLen := binary.LittleEndian.Uint16(headerBuf[5:7])
			pkgLen := uint16(headerBuf[3])
			if bodyLen > 0 {
				pkgLen += bodyLen + 2
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
				exportPacket := egtsschema.EgtsPackage{}
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
						readyToPersist = true
						logger.Debugf("Разбор подзаписи EGTS_SR_POS_DATA")
						exportPacket.MeasurementTimestamp = subRecData.NavigationTime.Unix()
						exportPacket.ReceivedTimestamp = time.Now().UTC().Unix()
						exportPacket.Latitude = subRecData.Latitude
						exportPacket.Longitude = subRecData.Longitude
						exportPacket.Speed = int32(subRecData.Speed)
						exportPacket.Direction = int32(subRecData.Direction)
						exportPacket.Guid = fmt.Sprintf("%s", uuid.New())
					case *egts.SrExtPosData:
						logger.Debugf("Разбор подзаписи EGTS_SR_EXT_POS_DATA")
						exportPacket.NumOfSatelites = int32(subRecData.Satellites)
						exportPacket.Pdop = int32(subRecData.PositionDilutionOfPrecision)
						exportPacket.Hdop = int32(subRecData.HorizontalDilutionOfPrecision)
						exportPacket.Vdop = int32(subRecData.VerticalDilutionOfPrecision)
						exportPacket.NavigationSystem = toNavigationSystem(subRecData.NavigationSystem)

					case *egts.SrAdSensorsData:
						readyToPersist = true
						logger.Debugf("Разбор подзаписи EGTS_SR_AD_SENSORS_DATA")
						analogSensors := egtsschema.UnionArrayAnalogSensorNull{}
						if subRecData.AnalogSensorFieldExists1 == "1" {
							sensor := egtsschema.AnalogSensor{1, int32(subRecData.AnalogSensor1)}
							analogSensors.ArrayAnalogSensor = append(analogSensors.ArrayAnalogSensor, &sensor)
						}
						if subRecData.AnalogSensorFieldExists2 == "1" {
							sensor := egtsschema.AnalogSensor{2, int32(subRecData.AnalogSensor2)}
							analogSensors.ArrayAnalogSensor = append(analogSensors.ArrayAnalogSensor, &sensor)
						}
						if subRecData.AnalogSensorFieldExists3 == "1" {
							sensor := egtsschema.AnalogSensor{3, int32(subRecData.AnalogSensor3)}
							analogSensors.ArrayAnalogSensor = append(analogSensors.ArrayAnalogSensor, &sensor)
						}
						if subRecData.AnalogSensorFieldExists4 == "1" {
							sensor := egtsschema.AnalogSensor{4, int32(subRecData.AnalogSensor4)}
							analogSensors.ArrayAnalogSensor = append(analogSensors.ArrayAnalogSensor, &sensor)
						}
						if subRecData.AnalogSensorFieldExists5 == "1" {
							sensor := egtsschema.AnalogSensor{5, int32(subRecData.AnalogSensor5)}
							analogSensors.ArrayAnalogSensor = append(analogSensors.ArrayAnalogSensor, &sensor)
						}
						if subRecData.AnalogSensorFieldExists6 == "1" {
							sensor := egtsschema.AnalogSensor{6, int32(subRecData.AnalogSensor6)}
							analogSensors.ArrayAnalogSensor = append(analogSensors.ArrayAnalogSensor, &sensor)
						}
						if subRecData.AnalogSensorFieldExists7 == "1" {
							sensor := egtsschema.AnalogSensor{7, int32(subRecData.AnalogSensor7)}
							analogSensors.ArrayAnalogSensor = append(analogSensors.ArrayAnalogSensor, &sensor)
						}
						if subRecData.AnalogSensorFieldExists8 == "1" {
							sensor := egtsschema.AnalogSensor{8, int32(subRecData.AnalogSensor8)}
							analogSensors.ArrayAnalogSensor = append(analogSensors.ArrayAnalogSensor, &sensor)
						}
						exportPacket.AnalogSensors = &analogSensors
					case *egts.SrAbsCntrData:
						readyToPersist = true
						logger.Debugf("Разбор подзаписи EGTS_SR_ABS_CNTR_DATA")

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
						readyToPersist = true
						logger.Debugf("Разбор подзаписи EGTS_SR_LIQUID_LEVEL_SENSOR")
						liquidSensors := egtsschema.UnionArrayLiquidSensorNull{}
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
						liquidSensors.ArrayLiquidSensor = append(liquidSensors.ArrayLiquidSensor, &sensor)
						exportPacket.LiquidSensors = &liquidSensors
					}
				}

				if readyToPersist {
					exportPacket.Imei = deviceImei
					verifyNullableAttributes(&exportPacket)
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
			_, _ = conn.Write(resp)

			logger.Debugf("Отправлен пакет EGTS_PT_RESPONSE: %X", resp)

			if len(srResultCodePkg) > 0 {
				_, _ = conn.Write(srResultCodePkg)
				logger.Debugf("Отправлен пакет EGTS_SR_RESULT_CODE: %X", resp)
			}
		case egts.PtResponsePacket:
			logger.Debug("Тип пакета EGTS_PT_RESPONSE")
		}
	}
}

func verifyNullableAttributes(egtsPackage *egtsschema.EgtsPackage) {
	if egtsPackage.AnalogSensors == nil {
		sensors := egtsschema.UnionArrayAnalogSensorNull{}
		egtsPackage.AnalogSensors = &sensors

	}
	if egtsPackage.LiquidSensors == nil {
		sensors := egtsschema.UnionArrayLiquidSensorNull{}
		egtsPackage.LiquidSensors = &sensors

	}
}

func toNavigationSystem(egtsNavSystemCode uint16) egtsschema.NavigationSystem {
	switch egtsNavSystemCode {
	// Glonass
	case 1:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemGLONASS)
	// GPS
	case 2:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemGPS)
	// Galileo
	case 4:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemGalileo)
	// Compass
	case 8:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemCompass)
	// Beidou
	case 16:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemBeidou)
	// DORIS
	case 32:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemDORIS)
	// unknown
	default: // including 0
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemUknown)
	}
}
