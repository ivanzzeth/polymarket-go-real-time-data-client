package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	polymarketdataclient "github.com/ivanzzeth/polymarket-go-real-time-data-client"
)

func main() {
	// Create a typed message router
	router := polymarketdataclient.NewRealtimeMessageRouter()

	// Track crypto prices
	router.RegisterCryptoPriceHandler(func(price polymarketdataclient.CryptoPrice) error {
		log.Printf("[Crypto] %s = $%s (time: %s)",
			price.Symbol,
			price.Value.String(),
			price.Time.Format("15:04:05.000"))
		return nil
	})

	// Track equity prices
	router.RegisterEquityPriceHandler(func(price polymarketdataclient.EquityPrice) error {
		log.Printf("[Equity] %s = $%s (time: %s)",
			price.Symbol,
			price.Value.String(),
			price.Time.Format("15:04:05.000"))
		return nil
	})

	// Create the WebSocket client
	client := polymarketdataclient.New(
		polymarketdataclient.WithLogger(polymarketdataclient.NewLogger()),
		polymarketdataclient.WithOnConnect(func() {
			log.Println("✓ Connected to Polymarket WebSocket!")
		}),
		polymarketdataclient.WithOnNewMessage(func(data []byte) {
			// Route the message to the appropriate typed handler
			if err := router.RouteMessage(data); err != nil {
				log.Printf("Error routing message: %v", err)
			}
		}),
	)

	// Connect to the server
	log.Println("Connecting to Polymarket...")
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	// Create typed subscription handler
	typedSub := polymarketdataclient.NewRealtimeTypedSubscriptionHandler(client)

	log.Println("Subscribing to price feeds...")

	// Subscribe to crypto prices using predefined constants
	cryptoFilters := []*polymarketdataclient.CryptoPriceFilter{
		polymarketdataclient.NewBTCPriceFilter(),
		polymarketdataclient.NewETHPriceFilter(),
		polymarketdataclient.NewSOLPriceFilter(),
	}
	for _, filter := range cryptoFilters {
		if err := typedSub.SubscribeToCryptoPrices(nil, filter); err != nil {
			log.Printf("Failed to subscribe to %s: %v", filter.Symbol, err)
		} else {
			log.Printf("✓ Subscribed to %s", filter.Symbol)
		}
	}

	// Subscribe to equity prices using predefined constants
	equityFilters := []*polymarketdataclient.EquityPriceFilter{
		polymarketdataclient.NewAppleStockFilter(),
		polymarketdataclient.NewTeslaStockFilter(),
		polymarketdataclient.NewNvidiaStockFilter(),
		polymarketdataclient.NewEquityPriceFilter(polymarketdataclient.EquitySymbolMSFT),
	}
	for _, filter := range equityFilters {
		if err := typedSub.SubscribeToEquityPrices(nil, filter); err != nil {
			log.Printf("Failed to subscribe to %s: %v", filter.Symbol, err)
		} else {
			log.Printf("✓ Subscribed to %s", filter.Symbol)
		}
	}

	log.Println("=== Price Tracking Started ===")
	log.Println("Monitoring crypto and equity prices...")
	log.Println("Press Ctrl+C to exit")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n\nShutting down...")
	log.Println("✓ Disconnected successfully")
}
