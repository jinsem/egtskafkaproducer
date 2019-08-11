package main

import (
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
)

type EgtsKafkaPersister struct{}

func (p *EgtsKafkaPersister) Initialize(kafka KafkaSettings) error {
	return nil
}

func (p *EgtsKafkaPersister) Produce(egtsPackage *egtsschema.EgtsPackage) error {
	return nil
}

func (p *EgtsKafkaPersister) Close() error {
	return nil
}
