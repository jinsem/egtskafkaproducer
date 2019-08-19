package main

import (
	"bytes"
	"context"
	"fmt"
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
	"github.com/landoop/schema-registry"
	"github.com/segmentio/kafka-go"
)

type EgtsKafkaPersister struct {
	kafkaSettings *KafkaSettings
	writer        *kafka.Writer
	// TODO: Do not cahche ID forever, but request it more frequently
	valueSchemaId uint32
}

func (p *EgtsKafkaPersister) Initialize(cfg *KafkaSettings) error {
	if cfg == nil {
		return fmt.Errorf("Configuration is not set")
	}
	p.kafkaSettings = cfg
	p.writer = p.getKafkaWriter(cfg)
	err := p.initSchemaId()
	if err == nil {
		logger.Debug("Kafka persister готов к работе")
	} else {
		logger.Error("Ошибка получения идентификаторов схемы")
	}
	return err
}

func (p *EgtsKafkaPersister) initSchemaId() error {
	client, _ := schemaregistry.NewClient(p.kafkaSettings.SchemaRegistryUrl)
	valueSubjectName := p.kafkaSettings.OutputTopicName + "-value"
	valueSchema, err := client.GetLatestSchema(valueSubjectName)
	if err == nil {
		p.valueSchemaId = uint32(valueSchema.ID)
	}
	return err
}

func (c *EgtsKafkaPersister) getKafkaWriter(cfg *KafkaSettings) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.OutputTopicName,
		Balancer: &kafka.LeastBytes{},
	})
}

func (p *EgtsKafkaPersister) Produce(egtsPackage *egtsschema.EgtsPackage) error {
	logger.Debug("Processing message... ")
	var buf bytes.Buffer
	addSchemaRegistryHeader(&buf, p.valueSchemaId)
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

func (p *EgtsKafkaPersister) Close() error {
	err := p.writer.Close()
	logger.Debug("Kafka connector stopped ")
	return err
}
