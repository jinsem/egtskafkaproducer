package main

type EgtsProducer struct{}

func (p *EgtsProducer) Initialize(kafka KafkaSettings) error {
	return nil
}

func (p *EgtsProducer) Produce() error {
	return nil
}

func (p *EgtsProducer) Close() error {
	return nil
}
