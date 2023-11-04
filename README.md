# homeqtt

I had the problem that values from my Aqara temperature sensors, connected to Home Assistant using zigbee2mqtt, 
disappeared after some time. Because of this, I decided to write this client to log all sensor data messages to
a SQLite database to allow long-term analysis and persistence of the data.

