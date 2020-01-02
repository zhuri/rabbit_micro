package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func main() {
	// Get the connection string from the environment var
	url := os.Getenv("AMQP_URL")

	//If it doesnt exist, use the default connection string
	if url == "" {
		//url = "amqp://guest:guest@0.0.0.0:5672"
    	url = "amqp://guest:guest@rabbitmq"

	}

	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	// Create a channel from the connection
	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}

	// ExchangeDeclare - exchange point
	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	// We create a queue named Test
	_, err = channel.QueueDeclare("test", true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("exit", text) == 0 {
			return
		}

		// We create a queue message of amqp publishing instance
		message := amqp.Publishing{
			Body: []byte(text),
		}

		// We publish the message to the EXCHANGE
		err = channel.Publish("events", "random-key", false, false, message)

		if err != nil {
			panic("error publishing a message to the queue:" + err.Error())
		}

		// We bind the queue to the exchange to send and receive data from the queue
		err = channel.QueueBind("test", "#", "events", false, nil)

		if err != nil {
			panic("error binding to the queue: " + err.Error())
		}
	}
}
