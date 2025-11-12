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

	// Register handler for equity prices
	router.RegisterEquityPriceHandler(func(price polymarketdataclient.EquityPrice) error {
		log.Printf("[Equity] %s = $%s (updated at %s)",
			price.Symbol,
			price.Value.String(),
			price.Time.Format("15:04:05"))
		return nil
	})

	// Create client with message callback
	client := polymarketdataclient.New(
		polymarketdataclient.WithLogger(polymarketdataclient.NewSilentLogger()),
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

	// Subscribe to crypto prices using predefined filters
	log.Println("\nSubscribing to crypto prices...")
	cryptoFilters := []*polymarketdataclient.CryptoPriceFilter{
		polymarketdataclient.NewBTCPriceFilter(),    // Bitcoin
		polymarketdataclient.NewETHPriceFilter(),    // Ethereum
		polymarketdataclient.NewSOLPriceFilter(),    // Solana
	}

	for _, filter := range cryptoFilters {
		if err := typedSub.SubscribeToCryptoPrices(nil, filter); err != nil {
			log.Printf("Failed to subscribe to %s: %v", filter.Symbol, err)
		} else {
			log.Printf("✓ Subscribed to %s", filter.Symbol)
		}
	}

	// Subscribe to equity prices using predefined filters
	log.Println("\nSubscribing to equity prices...")
	equityFilters := []*polymarketdataclient.EquityPriceFilter{
		polymarketdataclient.NewAppleStockFilter(),   // Apple
		polymarketdataclient.NewTeslaStockFilter(),   // Tesla
		polymarketdataclient.NewNvidiaStockFilter(),  // NVIDIA
	}

	for _, filter := range equityFilters {
		if err := typedSub.SubscribeToEquityPrices(nil, filter); err != nil {
			log.Printf("Failed to subscribe to %s: %v", filter.Symbol, err)
		} else {
			log.Printf("✓ Subscribed to %s", filter.Symbol)
		}
	}

	// You can also use custom symbols with constants
	log.Println("\nSubscribing to additional symbols using constants...")
	customCrypto := polymarketdataclient.NewCryptoPriceFilter(polymarketdataclient.CryptoSymbolAVAXUSDT)
	if err := typedSub.SubscribeToCryptoPrices(nil, customCrypto); err != nil {
		log.Printf("Failed to subscribe: %v", err)
	} else {
		log.Printf("✓ Subscribed to %s", customCrypto.Symbol)
	}

	log.Println("\n=== Now receiving live price updates ===")
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
