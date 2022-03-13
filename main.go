package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/segmentio/kafka-go"
)

const (
	topic = "send_email"
)

type Message struct {
	EmailSubject  string            `json:"subject"`
	EmailTemplate string            `json:"template"`
	EmailRcpt     []string          `json:"rcpt"`
	EmailData     map[string]string `json:"data"`
}

func main() {
	// create a new context
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		return
	}

	consume(ctx)
}

func consume(ctx context.Context) {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	l := log.New(os.Stdout, "kafka reader: ", 0)
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages

	brokerAddress := strings.Split(os.Getenv("BROKER_ADDRESS"), ",")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddress,
		Topic:   topic,
		GroupID: os.Getenv("MESSAGE_GROUP"),
		// assign the logger to the reader
		Logger: l,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		m := Message{}
		json.Unmarshal(msg.Value, &m)

		SendMail(m.EmailSubject, m.EmailTemplate, m.EmailRcpt, m.EmailData)

		// after receiving the message, log its value
		fmt.Println("received: ", m)
	}
}
