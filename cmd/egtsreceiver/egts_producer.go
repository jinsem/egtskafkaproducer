package main

import (
	egtsschema "../../pkg/avro"
)

type EgtsProducer struct{}

func (p *EgtsProducer) Initialize(kafka KafkaSettings) error {
	return nil
}

func (p *EgtsProducer) Produce(egtsPackage *egtsschema.EgtsPackage) error {
	return nil
}

func (p *EgtsProducer) Close() error {
	return nil
}
