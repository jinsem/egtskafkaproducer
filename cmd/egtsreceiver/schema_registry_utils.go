package main

import (
	"bytes"
)

func addSchemaRegistryHeader(buffer *bytes.Buffer, schemaId uint32) {
	// First byte must be 0
	buffer.WriteByte(0)
	buffer.WriteByte(byte(schemaId >> 24))
	buffer.WriteByte(byte(schemaId >> 16))
	buffer.WriteByte(byte(schemaId >> 8))
	buffer.WriteByte(byte(schemaId))
}
