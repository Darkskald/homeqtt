package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type HandlerBuilder func(sensorTopic string) mqtt.MessageHandler

func setUpMqttOptions(config Config, builder HandlerBuilder) *mqtt.ClientOptions {

	log.Println(config.MQTTEndpoint())

	opts := mqtt.NewClientOptions().AddBroker(config.MQTTEndpoint())

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("connection established")

		// Establish the subscription - doing this here means that it will happen every time a connection is established
		// (useful if opts.CleanSession is TRUE or the broker does not reliably store session data)

		for _, topic := range config.SplitSensorTopics() {

			fmt.Println(topic)

			t := c.Subscribe(topic, 1, builder(topic))
			// the connection handler is called in a goroutine so blocking here would hot cause an issue. However, as blocking
			// in other handlers does cause problems its best to just assume we should not block
			topic := topic
			go func() {
				_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
				if t.Error() != nil {
					log.Printf("ERROR SUBSCRIBING: %s\n", t.Error())
				} else {
					log.Println("subscribed to: ", topic)
				}
			}()
		}

	}
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		log.Println("attempting to reconnect")
	}

	return opts
}
