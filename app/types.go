package main

// Interface for EGTS kafka producer
type Producer interface {

	// Initialize producer
	Initialize(kafka KafkaSettings) error

	// Produce record
	// TODO: Define input parameter
	Produce() error

	// Close producer and all the opened resources
	Close() error
}
