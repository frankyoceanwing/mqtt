package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	hostname, _ := os.Hostname()
	server := flag.String("server", "tcp://127.0.0.1:1883", "The full URL of the MQTT server to connect to")
	topic := flag.String("topic", hostname, "Topic to publish the messages on")
	qos := flag.Int("qos", 1, "The QoS to send the messages at")
	retained := flag.Bool("retained", false, "Are the messages sent with the retained flag")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
	fmt.Printf("Connected to %s\n", *server)

	stdin := bufio.NewReader(os.Stdin)
	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		token := client.Publish(*topic, byte(*qos), *retained, message)
		fmt.Printf("publish [topic=\"%s\" message=\"%s\" complete=%t]\n",
			*topic, strings.Replace(message, "\n", "", -1), token.WaitTimeout(5*time.Second))
	}
}
