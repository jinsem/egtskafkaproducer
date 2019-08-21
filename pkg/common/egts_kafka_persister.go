package common

import (
	"bytes"
	"context"
	"fmt"
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
	"github.com/labstack/gommon/log"
	"github.com/landoop/schema-registry"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	kafkaSettings *KafkaSettings
	writer        *kafka.Writer
	valueSchemaId uint32
	logger        *log.Logger
}

func (p *KafkaProducer) Initialize(cfg *KafkaSettings, logger *log.Logger) error {
	if cfg == nil {
		return fmt.Errorf("Configuration is not set")
	}
	p.kafkaSettings = cfg
	p.writer = p.getKafkaWriter(cfg)
	err := p.initSchemaId()
	p.logger = logger
	if err == nil {
		p.logger.Debug("Kafka persister готов к работе")
	} else {
		p.logger.Error("Ошибка получения идентификатора схемы")
	}
	return err
}

func (p *KafkaProducer) initSchemaId() error {
	client, _ := schemaregistry.NewClient(p.kafkaSettings.SchemaRegistryUrl)
	valueSubjectName := p.kafkaSettings.OutputTopicName + "-value"
	schemaSource := egtsschema.EgtsPackage{}
	schemaId, err := RegisterSchemaIfNotExists(client, valueSubjectName, schemaSource.Schema())
	if err == nil {
		p.valueSchemaId = uint32(schemaId)
	}
	return err
}

func (c *KafkaProducer) getKafkaWriter(cfg *KafkaSettings) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.OutputTopicName,
		Balancer: &kafka.LeastBytes{},
	})
}

func (p *KafkaProducer) Produce(egtsPackage *egtsschema.EgtsPackage) error {
	p.logger.Debug("Processing message... ")
	var buf bytes.Buffer
	egtsPackage.Schema()
	AddSchemaRegistryHeader(&buf, p.valueSchemaId)
	err := egtsPackage.Serialize(&buf)
	if err == nil {
		innerPkg := buf.Bytes()
		kafkaMsg := kafka.Message{
			Key:   nil,
			Value: innerPkg,
		}
		err = p.writer.WriteMessages(context.Background(), kafkaMsg)
	}
	return err
}

func (p *KafkaProducer) Close() error {
	err := p.writer.Close()
	p.logger.Debug("Kafka connector stopped ")
	return err
}
