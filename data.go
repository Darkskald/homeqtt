package main

import "time"

type SensorData struct {
	Battery          int     `json:"battery"`
	Humidity         float64 `json:"humidity"`
	LinkQuality      int     `json:"linkquality"`
	PowerOutageCount int     `json:"power_outage_count"`
	Pressure         float64 `json:"pressure"`
	Temperature      float64 `json:"temperature"`
	Voltage          int     `json:"voltage"`
}

type WrappedData struct {
	SensorName string    `json:"sensor_name"`
	Timestamp  time.Time `json:"timestamp"`
	SensorData
}
