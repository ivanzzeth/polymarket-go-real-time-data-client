package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	polymarketdataclient "github.com/ivanzzeth/polymarket-go-real-time-data-client"
)

func main() {
	log.Println("=== Simple Typed Subscription Example ===")
	log.Println("This example shows how to use typed subscriptions with predefined symbols")
	log.Println()
	log.Println("IMPORTANT: crypto_prices and equity_prices topics only support ONE symbol per connection")
	log.Println("This example subscribes to a SINGLE symbol. See examples/multi_symbol_tracking for multiple symbols.")
	log.Println()

	// Symbol to monitor (you can change this)
	symbolToMonitor := polymarketdataclient.CryptoSymbolBTCUSDT // Options: CryptoSymbolBTCUSDT, CryptoSymbolETHUSDT, etc.

	// Create a message router for typed handling
	router := polymarketdataclient.NewRealtimeMessageRouter()

	// Register handler for crypto prices
	router.RegisterCryptoPriceHandler(func(price polymarketdataclient.CryptoPrice) error {
		log.Printf("[Crypto] %s = $%s (updated at %s)",
			price.Symbol,
			price.Value.String(),
			price.Time.Format("15:04:05"))
		return nil
	})

	// Create client with message callback
	client := polymarketdataclient.New(
		polymarketdataclient.WithLogger(polymarketdataclient.NewLogger()),
		polymarketdataclient.WithAutoReconnect(true),
		polymarketdataclient.WithOnConnect(func() {
			log.Println("✅ Connected to Polymarket WebSocket")
		}),
		polymarketdataclient.WithOnNewMessage(func(data []byte) {
			router.RouteMessage(data)
		}),
	)

	// Connect to the server
	log.Println("Connecting to Polymarket WebSocket...")
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Create typed subscription handler
	typedSub := polymarketdataclient.NewRealtimeTypedSubscriptionHandler(client)

	// Subscribe to a single crypto price using predefined filter
	log.Printf("\nSubscribing to %s...\n", symbolToMonitor)
	filter := polymarketdataclient.NewCryptoPriceFilter(symbolToMonitor)
	if err := typedSub.SubscribeToCryptoPrices(nil, filter); err != nil {
		log.Fatalf("Failed to subscribe to %s: %v", symbolToMonitor, err)
	}
	log.Printf("✓ Successfully subscribed to %s\n", symbolToMonitor)

	log.Println()
	log.Println("=== Now receiving live price updates ===")
	log.Printf("Monitoring: %s\n", symbolToMonitor)
	log.Println("Press Ctrl+C to exit")
	log.Println()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\nShutting down...")
	if err := client.Disconnect(); err != nil {
		log.Printf("Error during disconnect: %v", err)
	}
	log.Println("✅ Disconnected successfully")
}
