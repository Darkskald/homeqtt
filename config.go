package main

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
	"strings"
)

type Config struct {
	MqttHost     string `env:"MQTTHOST,required"`
	MqttPort     string `env:"MQTTPORT,required"`
	SensorTopics string `env:"SENSORTOPICS,required"`
	DatabasePath string `env:"DBPATH,required"`
}

func ParseConfig() Config {
	godotenv.Load()
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	return cfg
}

func (c Config) MQTTEndpoint() string {
	return fmt.Sprintf("tcp://%s:%s", c.MqttHost, c.MqttPort)
}

func (c Config) SplitSensorTopics() []string {
	return strings.Split(c.SensorTopics, ",")
}

func GetSensorName(topic string) string {
	// TODO add checks here
	return strings.Split(topic, "/")[1]
}
