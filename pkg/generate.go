package pkg

// Use a go:generate directive to build the Go structs for `measurementpackage.avsc`

//go:generate mkdir -p ./avro
//go:generate $GOPATH/bin/gogen-avro ./avro schema/measurementpackage.avsc
