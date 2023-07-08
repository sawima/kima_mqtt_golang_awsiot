package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sawima/kima_mqtt_golang_awsiot/imaiot"
)

func main() {
	imaiot.ConnectIOT()

	quit := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("[MQTT] Disconnected")
		quit <- struct{}{}
	}()
	<-quit
}
