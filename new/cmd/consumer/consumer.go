package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/streadway/amqp"
)

type AMQ struct {
	amqpServer      string
	amqpVirtualHost string
	amqpUsername    string
	amqpPassword    string
	amqpQueue       string
}

func (a *AMQ) StringConnection() string {
	return "amqp://" + a.amqpUsername + ":" + a.amqpPassword + "@" + a.amqpServer + ":5672/" + a.amqpVirtualHost
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	amqpConnect *AMQ
)

func init() {
	amqpConnect = &AMQ{
		amqpServer:      os.Getenv("AMQP_SERVER"),
		amqpVirtualHost: os.Getenv("AMQP_VHOST"),
		amqpUsername:    os.Getenv("AMQP_USERNAME"),
		amqpPassword:    os.Getenv("AMQP_PASSWORD"),
		amqpQueue:       os.Getenv("AMQP_QUEUE"),
	}
}

func main() {
	conn, err := amqp.Dial(amqpConnect.StringConnection())
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		amqpConnect.amqpQueue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
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

	forever := make(chan bool)

	go func() {
		wg := new(sync.WaitGroup)
		for d := range msgs {
			log.Printf("Received a message: %s", d.MessageId)
			wg.Add(1)
			execProccesser(d.Body, wg)
		}
		wg.Wait()
	}()

	log.Printf(" [*] Waiting for messages from %s", amqpConnect.amqpVirtualHost)
	<-forever
}

func execProccesser(msg []byte, wg *sync.WaitGroup) {
	tempDir, err := ioutil.TempDir("/tmp", "tmp.")
	failOnError(err, "Failed to create temporal directory")

	mFile := tempDir + "/message.msg"
	err = ioutil.WriteFile(mFile, msg, 0400)
	failOnError(err, "Failed to write message file")

	cmd := exec.Command("/root/proccesser.sh", mFile)
	log.Printf("Command to execute: %s", cmd.Args)
	out, err := cmd.Output()
	fmt.Printf("%s", out)
	log.Printf("Command finished with error: %v", err)
	wg.Done()
}
