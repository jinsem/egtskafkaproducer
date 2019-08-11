package main

// Use a go:generate directive to build the Go structs for `egtsmessage.avsc`

//go:generate mkdir -p ./avro
//go:generate $GOPATH/bin/gogen-avro ./avro schema/egtsmessage.avsc
