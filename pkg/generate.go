package pkg

// Use a go:generate directive to build the Go structs for `egtspackage.avsc`

//go:generate mkdir -p ./avro
//go:generate $GOPATH/bin/gogen-avro ./avro schema/egtspackage.avsc
