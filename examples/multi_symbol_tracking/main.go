package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	polymarketdataclient "github.com/ivanzzeth/polymarket-go-real-time-data-client"
)

func main() {
	log.Println("=== Multi-Symbol Price Tracking Demo ===")
	log.Println()
	log.Println("IMPORTANT: Due to Polymarket API limitations, crypto_prices and equity_prices")
	log.Println("topics only support ONE symbol per WebSocket connection.")
	log.Println()
	log.Println("This example demonstrates how to monitor MULTIPLE symbols by creating")
	log.Println("separate WebSocket connections for each symbol.")
	log.Println()

	// Symbols to monitor
	symbols := []string{"btcusdt", "ethusdt", "solusdt"}

	// Shared message counter
	var messageCount sync.Map // map[string]int

	// Create separate client for each symbol
	var clients []polymarketdataclient.WsClient
	var wg sync.WaitGroup

	for _, symbol := range symbols {
		symbol := symbol // capture loop variable
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Create a message router for this symbol
			router := polymarketdataclient.NewRealtimeMessageRouter()

			// Register handler for this symbol
			router.RegisterCryptoPriceHandler(func(price polymarketdataclient.CryptoPrice) error {
				// Increment message count
				count := 1
				if val, ok := messageCount.Load(symbol); ok {
					count = val.(int) + 1
				}
				messageCount.Store(symbol, count)

				log.Printf("[%s] Price: $%s (update #%d)",
					price.Symbol,
					price.Value.String(),
					count)
				return nil
			})

			// Create dedicated client for this symbol
			client := polymarketdataclient.New(
				// polymarketdataclient.WithLogger(polymarketdataclient.NewLogger()),
				polymarketdataclient.WithOnConnect(func() {
					log.Printf("✓ [%s] Connected", symbol)
				}),
				polymarketdataclient.WithOnDisconnect(func(err error) {
					if err != nil {
						log.Printf("✗ [%s] Disconnected: %v", symbol, err)
					}
				}),
				polymarketdataclient.WithOnNewMessage(func(data []byte) {
					if err := router.RouteMessage(data); err != nil {
						log.Printf("⚠️ [%s] Error routing message: %v", symbol, err)
					}
				}),
			)

			// Connect to the server
			log.Printf("Connecting [%s]...", symbol)
			if err := client.Connect(); err != nil {
				log.Printf("Failed to connect [%s]: %v", symbol, err)
				return
			}

			// Store client for cleanup
			clients = append(clients, client)

			// Create typed subscription handler
			typedSub := polymarketdataclient.NewRealtimeTypedSubscriptionHandler(client)

			// Subscribe to this specific symbol
			filter := polymarketdataclient.NewCryptoPriceFilter(symbol)
			if err := typedSub.SubscribeToCryptoPrices(nil, filter); err != nil {
				log.Printf("Failed to subscribe to [%s]: %v", symbol, err)
				client.Disconnect()
				return
			}
			log.Printf("✓ [%s] Subscribed successfully", symbol)
		}()
	}

	// Wait for all connections to be established
	wg.Wait()

	log.Println()
	log.Println("=== Multi-Symbol Tracking Started ===")
	log.Printf("Monitoring %d symbols simultaneously: %v\n", len(symbols), symbols)
	log.Println("Each symbol has its own WebSocket connection.")
	log.Println("Press Ctrl+C to exit")
	log.Println()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println()
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("Shutting down...")
	log.Println()

	// Print final statistics
	log.Println("📊 Final Statistics:")
	for _, symbol := range symbols {
		if val, ok := messageCount.Load(symbol); ok {
			log.Printf("  %s: %d updates received", symbol, val.(int))
		} else {
			log.Printf("  %s: 0 updates received", symbol)
		}
	}
	log.Println()

	// Disconnect all clients
	for i, client := range clients {
		if err := client.Disconnect(); err != nil {
			log.Printf("Error disconnecting client %d: %v", i, err)
		}
	}

	log.Println("✓ All connections closed successfully")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
