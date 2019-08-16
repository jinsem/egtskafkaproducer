#!/bin/bash

docker run \
    --rm \
    --net=host \
    --name=egtsreceiver \
    -e EGTS_RECEIVER_APP_HOSTNAME= \
    -e EGTS_RECEIVER_APP_PORT=6000 \
    -e EGTS_RECEIVER_APP_CONNECTIONTIMETOLIVESEC=60 \
    -e EGTS_RECEIVER_LOG_LEVEL=DEBUG \
    -e EGTS_RECEIVER_KAFKA_BROKERS=localhost:29092 \
    -e EGTS_RECEIVER_KAFKA_OUTPUTTOPICNAME=sensor_data \
    testegts:9

# esemionov/egtskafkaproducer:latest
