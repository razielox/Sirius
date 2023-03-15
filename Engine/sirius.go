// Sirius Scanning Engine
package main

import (
	"encoding/json"
	"fmt"
	"log"

	core "github.com/0sm0s1z/Sirius-Scan/Engine/core"
	lib "github.com/0sm0s1z/Sirius-Scan/Engine/lib"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Engine is the main scanning engine
func main() {
	fmt.Println("Sirius Scanning Engine")

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"scan", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	/*
		Main loop for the engine
			- Listen for messages
			- Execute scan based on message contents
			- Send results back to the queue
	*/
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Execute & Manage scans based on massage contents
			var scanRequest lib.ScanRequest
			err := json.Unmarshal(d.Body, &scanRequest)
			if err != nil {
				fmt.Println("JSON Unmarshal format error!", err)
			}

			// Case Statement for scan tracking
			switch scanRequest.Command {
			case "new":
				// Execute the scan
				log.Println("=== New Scan Requested ===")
				// Spawn a new scan
				go func() {
					core.NewScan(scanRequest)
				}()
			case "report":
				// Stop the scan
				log.Println("Scan Stopped")
			case "scanDiscovery":
				// Stop the scan
				log.Println("Scan Stopped")
			case "scanVulnerability":
				// Stop the scan
				log.Println("Scan Stopped")
			}
		}
	}()

	log.Printf(" [*] Sirius Scanning Engine Started. Waiting for scan requests. To exit press CTRL+C")
	<-forever
}
