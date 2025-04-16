package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	uri           string
	exchange      string
	routingKey    string
	body          string
	inputFilePath string
	persistent    bool // New flag for delivery mode
)

const (
	NonPersistent uint8 = 1
	Persistent    uint8 = 2
)

func init() {
	flag.StringVar(&uri, "uri", "", "AMQP URI amqp://<user>:<password>@<host>:<port>/[vhost]")
	flag.StringVar(&exchange, "exchange", "", "Exchange name")
	flag.StringVar(&routingKey, "routing-key", "", `Routing key. Use queue name with blank exchange to publish directly to queue. Use queue name with blank exchange to publish directly to queue.`)
	flag.StringVar(&body, "body", "", "Message body")
	flag.StringVar(&inputFilePath, "input-file", "", "Input file path")
	flag.BoolVar(&persistent, "persistent", false, "Use persistent delivery mode (default: non-persistent)") // Add the new flag

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "A command-line tool for publishing messages to a RabbitMQ server.")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := validateFlags(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validateFlags() error {
	if uri == "" {
		return errors.New("uri cannot be blank")
	}
	if exchange == "" && routingKey == "" {
		return errors.New("exchange and routing-key cannot both be blank")
	}
	if body == "" && inputFilePath == "" {
		return errors.New("body and input-file cannot both be blank")
	}
	return nil
}

func getMessages() ([]string, error) {
	if inputFilePath != "" {
		content, err := os.ReadFile(inputFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read input file: %w", err)
		}
		lines := strings.Split(string(content), "\n")
		messages := make([]string, 0, len(lines)) // Pre-allocate slice
		for _, line := range lines {
			line = strings.TrimSpace(line) // Remove leading/trailing whitespace
			if line != "" {
				messages = append(messages, line)
			}
		}
		return messages, nil
	}
	return []string{body}, nil
}

func main() {
	connection, err := amqp.Dial(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer channel.Close() // close the channel as well

	messages, err := getMessages()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Printf("%d messages to publish", len(messages))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Adjust timeout as needed
	defer cancel()

	deliveryMode := NonPersistent // Default to non-persistent
	if persistent {
		deliveryMode = Persistent
	}

	for _, m := range messages {
		err = channel.PublishWithContext(
			ctx,
			exchange,
			routingKey,
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				DeliveryMode: deliveryMode, // Set the delivery mode
				ContentType:  "text/plain",
				Body:         []byte(m),
			})
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			// Consider whether to continue publishing or exit
		}
	}

	log.Println("Messages published successfully.")
}
