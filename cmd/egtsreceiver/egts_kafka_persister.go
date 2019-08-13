package main

import (
	"bytes"
	"context"
	"fmt"
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
	"github.com/segmentio/kafka-go"
)

type EgtsKafkaPersister struct {
	kafkaSettings *KafkaSettings
	writer        *kafka.Writer
}

func (p *EgtsKafkaPersister) Initialize(cfg *KafkaSettings) error {
	if cfg == nil {
		return fmt.Errorf("Configuration is not set")
	}
	p.kafkaSettings = cfg
	p.writer = p.getKafkaWriter(cfg)
	logger.Debug("Kafka connector is ready")
	return nil
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
