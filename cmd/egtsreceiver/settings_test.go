package main

import (
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfLoad(t *testing.T) {
	configuration :=
		`app:
  HostName: "localhost:9092"
  Port: 6000
  ConnectionTimeToLiveSec: 60
log:
  Level: "DEBUG"
kafka:
  Brokers:
    - "localhost1:9092"
    - "localhost2:9092"
  OutputTopicName: "egts_package_sample"`
	file, err := ioutil.TempFile("/tmp", "config-for-test*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if _, err = file.WriteString(configuration); err != nil {
		t.Fatal(err)
	}
	settings := Settings{}
	if err = settings.LoadFromFile(file.Name()); err != nil {
		t.Fatal(err)
	}

	compareToSettings := Settings{
		App: AppSettings{
			HostName:                "localhost:9092",
			Port:                    "6000",
			ConnectionTimeToLiveSec: 60,
		},
		Log: LogSettings{
			Level: "DEBUG",
		},
		Kafka: KafkaSettings{
			Brokers:         []string{"localhost1:9092", "localhost2:9092"},
			OutputTopicName: "egts_package_sample",
		},
	}
	if diff := cmp.Diff(compareToSettings, settings); diff != "" {
		t.Errorf("Loaded configuration content is not expected: (-expected +current)\n%s", diff)
	}
}

func TestConfNotFound(t *testing.T) {
	_ = os.Setenv("EGTS_RECEIVER_APP_HOSTNAME", "localhost:9092")
	_ = os.Setenv("EGTS_RECEIVER_APP_PORT", "6000")
	_ = os.Setenv("EGTS_RECEIVER_APP_CONNECTIONTIMETOLIVESEC", "60")
	_ = os.Setenv("EGTS_RECEIVER_LOG_LEVEL", "DEBUG")
	_ = os.Setenv("EGTS_RECEIVER_KAFKA_BROKERS", "localhost1:9092,localhost2:9092")
	_ = os.Setenv("EGTS_RECEIVER_KAFKA_OUTPUTTOPICNAME", "egts_package_sample")

	settings := Settings{}

	if err := settings.LoadFromEnv(); err != nil {
		t.Fatal(err)
	}

	compareToSettings := Settings{
		App: AppSettings{
			HostName:                "localhost:9092",
			Port:                    "6000",
			ConnectionTimeToLiveSec: 60,
		},
		Log: LogSettings{
			Level: "DEBUG",
		},
		Kafka: KafkaSettings{
			Brokers:         []string{"localhost1:9092", "localhost2:9092"},
			OutputTopicName: "egts_package_sample",
		},
	}
	if diff := cmp.Diff(compareToSettings, settings); diff != "" {
		t.Errorf("Loaded configuration content is not expected: (-expected +current)\n%s", diff)
	}
}
