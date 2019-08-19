#!/usr/bin/env bash
#/bin/bash

BUILD_DEST=./bin

mkdir -p $BUILD_DEST
go build -o $BUILD_DEST/egtsreceiver ../cmd/egtsreceiver
