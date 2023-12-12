package main

import (
	"context"
	"flag"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	bodyFlag = flag.String("body", "Hello World!", "Message body to send")
	iterFlag = flag.Int("i", 1, "Number of send repetition")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	flag.Parse()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// body := "Hello World!"
	for i := 0; i < *iterFlag; i++ {
		err = ch.PublishWithContext(
			ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(*bodyFlag),
			},
		)
		log.Printf("%d. Send message: %s", i, *bodyFlag)
	}
	failOnError(err, "failed to publish")
}
