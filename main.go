package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func handler(message map[string]interface{}) {

	// decode the map values to string variables
	id := message["id"].(string)
	state := message["state"].(string)

	url := fmt.Sprintf("http://127.0.0.1:8697/api/vms/%s/power", id)
	method := "PUT"

	payload := strings.NewReader(state)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	authString := fmt.Sprintf("Basic %s", viper.Get("AUTH_TOKEN").(string))
	req.Header.Add("Accept", "application/vnd.vmware.vmw.rest-v1+json")
	req.Header.Add("Content-Type", "application/vnd.vmware.vmw.rest-v1+json")
	req.Header.Add("Authorization", authString)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

var knt int
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	messageData := msg.Payload()

	// expect message as in
	// data = '{"id": "ABCDEFGHIJKLMNOPQRSTUVWXYZ", state: "on"}'
	var message map[string]interface{}
	err := json.Unmarshal(messageData, &message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Bulding request to handler with id: '%s' and state: '%s'\n", message["id"], message["state"])

	// pass the map to the handler
	handler(message)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {

	// load config
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// start mqtt listener
	knt = 0
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	connectionString := fmt.Sprintf("tcp://%s:%s", viper.Get("MQTT_SERVER").(string), viper.Get("MQTT_PORT").(string))
	opts := mqtt.NewClientOptions().AddBroker(connectionString)
	opts.SetClientID("go_client")
	opts.SetUsername(viper.Get("MQTT_USERNAME").(string))
	opts.SetPassword(viper.Get("MQTT_PASSWORD").(string))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnectionLost = connectLostHandler
	topic := viper.Get("MQTT_TOPIC").(string)

	opts.OnConnect = func(c mqtt.Client) {
		if token := c.Subscribe(topic, 0, messagePubHandler); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}
	<-c
}
