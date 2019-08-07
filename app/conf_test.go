package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestConfLoad(t *testing.T) {
	configuration := `
app:
  host_name: ""
  port: 6000
  connection_ttl_sec: 60
  log:
    level: "DEBUG"
kafka:
  brokers:
    - "localhost:9092"
  output_topic_name: "egts_package_sample"
`
	file, err := ioutil.TempFile("/tmp", "config-for-test.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if _, err = file.WriteString(configuration); err != nil {
		t.Fatal(err)
	}
	settings := appSettings{}
	if settings.HostName == "" {

	}
	if err = settings.Load(file.Name()); err != nil {
		t.Fatal(err)
	}

}
