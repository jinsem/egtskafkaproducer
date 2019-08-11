package main

import (
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
)

// Interface for EGTS records persister
type Persister interface {
	Initialize(kafka KafkaSettings) error

	Persist(egtsPackage egtsschema.EgtsPackage) error

	Close() error
}
