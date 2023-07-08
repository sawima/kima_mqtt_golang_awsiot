package imaiot

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	//go:embed certs/AmazonRootCA1.pem
	pemCerts []byte
	//go:embed certs/e6be0a29bc-certificate.pem.crt
	certPem []byte
	//go:embed certs/e6be0a29bc-private.pem.key
	certKey     []byte
	cId         string //terminal ID
	stateTopic  string
	actionTopic string
	accessPoint string = "tls://x*x.iot.cn-northwest-1.amazonaws.com.cn:8883"
)

func init() {
	cId = "customizdClientID"
	stateTopic = "state/" + cId
	actionTopic = "thing/" + cId + "/action"
}

func ConnectIOT() error {
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(pemCerts)
	cer, err := tls.X509KeyPair(certPem, certKey)
	if err != nil {
		return err
	}

	willData := map[string]interface{}{}
	willData["state"] = "offline"
	willLoad, _ := json.Marshal(willData)

	config := &tls.Config{
		RootCAs:      certpool,
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    nil,
		Certificates: []tls.Certificate{cer},
	}

	connOpts := MQTT.NewClientOptions()
	if len(cId) == 0 {
		return errors.New("terminal is not correctly initialized")
	}
	connOpts.AddBroker(accessPoint)
	connOpts.SetClientID(cId).SetTLSConfig(config)
	connOpts.SetDefaultPublishHandler(pubHandler)
	connOpts.SetWill(stateTopic, string(willLoad), 0, true)
	connOpts.SetCleanSession(true)
	connOpts.SetMaxReconnectInterval(10 * time.Second)
	connOpts.SetKeepAlive(60)
	connOpts.SetOnConnectHandler(iotOnConnect)
	connOpts.SetConnectionLostHandler(disconnectFunc)
	iotC := MQTT.NewClient(connOpts)
	if token := iotC.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
		return token.Error()
	}

	if token := iotC.Subscribe(actionTopic, 0, deviceActionHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}

	if token := iotC.Subscribe("state/123", 0, pubHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}

	fmt.Println("connected and waiting mqtt msg.")
	return nil
}

var pubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var iotOnConnect MQTT.OnConnectHandler = func(c MQTT.Client) {
	log.Println("setup connection")
	stateData, _ := json.Marshal(map[string]string{
		"state": "online",
	})

	c.Publish(stateTopic, 0, true, stateData)
}

var disconnectFunc MQTT.ConnectionLostHandler = func(c MQTT.Client, err error) {
	stateData, _ := json.Marshal(map[string]string{
		"state": "offline",
	})
	c.Publish(stateTopic, 0, true, stateData)
	log.Println("client is disconnected from aws")
}

var deviceActionHandler MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	fmt.Printf("action topic: %s\n", m.Payload())
}
