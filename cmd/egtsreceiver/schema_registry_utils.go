package main

import (
	"bytes"
	"github.com/landoop/schema-registry"
)

// AddSchemaRegistryHeader puts header to the message to make it compatible with Kafka avro serde that
// uses schema registry. The header consists of 5 bytes:
// 1st byte is always zero.
// bytes from 2 to 5 contains ID of the schema defined in schema registry
func AddSchemaRegistryHeader(buffer *bytes.Buffer, schemaId uint32) {
	// First byte must be 0
	buffer.WriteByte(0)
	buffer.WriteByte(byte(schemaId >> 24))
	buffer.WriteByte(byte(schemaId >> 16))
	buffer.WriteByte(byte(schemaId >> 8))
	buffer.WriteByte(byte(schemaId))
}

func RegisterSchemaIfNotExists(schemaRegistryClient *schemaregistry.Client, subject string, schema string) (uint32, error) {
	registered, existedSchema, err := schemaRegistryClient.IsRegistered(subject, schema)
	if err != nil {
		return uint32(0), err
	}
	if registered {
		return uint32(existedSchema.ID), nil
	}
	newSchemaId, err := schemaRegistryClient.RegisterNewSchema(subject, schema)
	if err == nil {
		return uint32(newSchemaId), nil
	} else {
		return uint32(0), err
	}
}
