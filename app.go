package main

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	db     *pgxpool.Pool
	client mqtt.Client
}

func NewApp() *App {

	app := App{}

	cfg := ParseConfig()
	opts := setUpMqttOptions(cfg, app.buildHandler)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	app.client = client

	conn, err := pgxpool.New(context.Background(), cfg.ProvidePostgresUrl())
	if err != nil {
		log.Fatal(err)
	}

	app.db = conn

	return &app
}

func (a *App) Run() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	log.Println("signal caught - exiting")
	a.client.Disconnect(1000)
	defer a.db.Close()

	log.Println("shutdown complete")
}

func (a *App) buildHandler(sensorTopic string) mqtt.MessageHandler {
	sensorName := GetSensorName(sensorTopic)
	log.Println("got sensor name ", sensorName)

	return func(_ mqtt.Client, msg mqtt.Message) {
		var sensorData SensorData
		if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
			log.Printf("message %s is not a valid event JSON: %s", msg.Payload(), err)
			return
		}
		go func() {
			wrappedData := WrappedData{
				SensorName: sensorName,
				Timestamp:  time.Now(),
				SensorData: sensorData,
			}

			log.Printf("got %+v \n", wrappedData)
			err := a.persistDatapoint(wrappedData)

			if err != nil {
				log.Printf("error while persisting %+v: %s", wrappedData, err.Error())
			}
		}()

	}
}

func (a *App) persistDatapoint(data WrappedData) error {
	_, err := a.db.Exec(context.Background(), `INSERT INTO sensor_data (SensorName, Timestamp, Battery, Humidity, LinkQuality, 
                         PowerOutageCount, Pressure, Temperature, Voltage)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, data.SensorName, data.Timestamp, data.Battery,
		data.Humidity, data.LinkQuality, data.PowerOutageCount, data.Pressure, data.Temperature, data.Voltage)
	if err != nil {
		return err
	}
	return nil
}
