package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	polymarketdataclient "github.com/ivanzzeth/polymarket-go-real-time-data-client"
)

func main() {
	// Create a new client with options
	client := polymarketdataclient.New(
		// polymarketdataclient.WithLogger(polymarketdataclient.NewLogger()),
		polymarketdataclient.WithLogger(polymarketdataclient.NewSilentLogger()),
		polymarketdataclient.WithOnConnect(func() {
			log.Println("Connected to Polymarket WebSocket!")
		}),
		polymarketdataclient.WithOnNewMessage(func(data []byte) {
			// log.Printf("Received raw message: %s\n", string(data))

			var msg polymarketdataclient.SubscriptionMessage
			err := json.Unmarshal(data, &msg)
			if err != nil {
				log.Printf("Invalid message %v received: %v", msg, err)
				return
			}

			// fmt.Printf("Received msg: %s\n", string(data))

			switch msg.Topic {
			case polymarketdataclient.TopicActivity:
				switch msg.Type {
				case polymarketdataclient.MessageTypeTrades:
					var trade polymarketdataclient.Trade
					err = json.Unmarshal(msg.Payload, &trade)
					if err != nil {
						log.Printf("Invalid trade %v received: %v", msg.Payload, err)
						return
					}

					log.Printf("Trade: %+v\n", trade)
				}
			case polymarketdataclient.TopicComments:
				// TODO:
			}

			// Handle any further message processing. Can use the types in `payload.go` to unmarshal
		}),
	)

	// Connect to the server
	if err := client.Connect(); err != nil {
		panic(err)
	}

	// market you want to filter
	event_slug := "btc-updown-15m-1762929900"

	// Subscribe to market data
	subscriptions := []polymarketdataclient.Subscription{
		{
			Topic: polymarketdataclient.TopicActivity,
			// Type:  polymarketdataclient.MessageTypeAll,
			Type:    polymarketdataclient.MessageTypeTrades,
			Filters: fmt.Sprintf(`{"event_slug":"%s"}`, event_slug),
		},
		{
			Topic: polymarketdataclient.TopicComments,
			Type:  polymarketdataclient.MessageTypeCommentCreated,
		},
	}

	if err := client.Subscribe(subscriptions); err != nil {
		panic(err)
	}

	// Keep the connection alive
	time.Sleep(30 * time.Second)

	// Clean up
	client.Disconnect()
}
