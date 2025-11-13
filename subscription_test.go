// +build integration

package polymarketrealtime

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// TestRawSubscriptionConcepts tests different subscription approaches at the raw WebSocket level
//
// FINDINGS FROM THESE TESTS:
// The Polymarket WebSocket API for crypto_prices, crypto_prices_chainlink, and equity_prices
// topics has a critical limitation: it only supports ONE symbol per WebSocket connection.
//
// Test Results:
// 1. ✗ Single message with multiple subscriptions → Only last symbol receives updates
// 2. ✗ Multiple messages with one subscription each → Only last symbol receives updates
// 3. ✗ Single subscription with multiple symbols in filter → Server rejects the connection
//
// SOLUTION: To monitor multiple symbols, you must create separate WebSocket connections,
// one for each symbol you want to track.
func TestRawSubscriptionConcepts(t *testing.T) {
	t.Run("SingleMessageMultipleSubscriptions", func(t *testing.T) {
		t.Log("=== Test: Send multiple subscriptions in ONE message ===")

		var receivedSOL, receivedETH, receivedBTC bool
		var mu sync.Mutex

		// Connect to WebSocket
		conn, _, err := websocket.DefaultDialer.Dial("wss://ws-live-data.polymarket.com", nil)
		require.NoError(t, err)
		defer conn.Close()

		t.Log("✓ Connected")

		// Send ONE subscription message with MULTIPLE subscriptions
		subscribeMsg := map[string]interface{}{
			"action": "subscribe",
			"subscriptions": []map[string]string{
				{
					"topic":   "crypto_prices",
					"type":    "update",
					"filters": `{"symbol":"solusdt"}`,
				},
				{
					"topic":   "crypto_prices",
					"type":    "update",
					"filters": `{"symbol":"ethusdt"}`,
				},
				{
					"topic":   "crypto_prices",
					"type":    "update",
					"filters": `{"symbol":"btcusdt"}`,
				},
			},
		}

		msgBytes, _ := json.Marshal(subscribeMsg)
		t.Logf("Sending: %s", string(msgBytes))

		err = conn.WriteMessage(websocket.TextMessage, msgBytes)
		require.NoError(t, err)

		t.Log("Waiting for responses...")

		// Read messages for 15 seconds
		done := make(chan struct{})
		go func() {
			time.Sleep(15 * time.Second)
			close(done)
		}()

	readLoop:
		for {
			select {
			case <-done:
				break readLoop
			default:
				conn.SetReadDeadline(time.Now().Add(2 * time.Second))
				msgType, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						t.Logf("WebSocket error: %v", err)
						break readLoop
					}
					// Timeout or other errors, continue
					continue
				}

				if msgType != websocket.TextMessage {
					continue
				}

				var msg map[string]interface{}
				if err := json.Unmarshal(message, &msg); err != nil {
					continue
				}

				// Check if it's an update message
				if msgType, ok := msg["type"].(string); ok && msgType == "update" {
					if payload, ok := msg["payload"].(map[string]interface{}); ok {
						if symbol, ok := payload["symbol"].(string); ok {
							mu.Lock()
							switch symbol {
							case "solusdt":
								if !receivedSOL {
									receivedSOL = true
									t.Logf("✓ Received SOL update")
								}
							case "ethusdt":
								if !receivedETH {
									receivedETH = true
									t.Logf("✓ Received ETH update")
								}
							case "btcusdt":
								if !receivedBTC {
									receivedBTC = true
									t.Logf("✓ Received BTC update")
								}
							}
							mu.Unlock()
						}
					}
				}
			}
		}
		mu.Lock()
		gotSOL := receivedSOL
		gotETH := receivedETH
		gotBTC := receivedBTC
		mu.Unlock()

		t.Logf("\nResults: SOL=%v, ETH=%v, BTC=%v", gotSOL, gotETH, gotBTC)

		// Assert: Only the LAST subscription (BTC) should receive updates
		// This proves that multiple subscriptions in one message don't work
		require.False(t, gotSOL, "SOL should NOT receive updates (only last subscription works)")
		require.False(t, gotETH, "ETH should NOT receive updates (only last subscription works)")
		require.True(t, gotBTC, "BTC should receive updates (it's the last subscription)")
	})

	t.Run("MultipleMessagesEachWithOneSubscription", func(t *testing.T) {
		t.Log("=== Test: Send THREE separate messages, each with ONE subscription ===")

		var receivedSOL, receivedETH, receivedBTC bool
		var mu sync.Mutex

		// Connect to WebSocket
		conn, _, err := websocket.DefaultDialer.Dial("wss://ws-live-data.polymarket.com", nil)
		require.NoError(t, err)
		defer conn.Close()

		t.Log("✓ Connected")

		// Send THREE separate subscription messages
		subscriptions := []map[string]string{
			{
				"topic":   "crypto_prices",
				"type":    "update",
				"filters": `{"symbol":"solusdt"}`,
			},
			{
				"topic":   "crypto_prices",
				"type":    "update",
				"filters": `{"symbol":"ethusdt"}`,
			},
			{
				"topic":   "crypto_prices",
				"type":    "update",
				"filters": `{"symbol":"btcusdt"}`,
			},
		}

		for i, sub := range subscriptions {
			subscribeMsg := map[string]interface{}{
				"action":        "subscribe",
				"subscriptions": []map[string]string{sub},
			}

			msgBytes, _ := json.Marshal(subscribeMsg)
			t.Logf("Message %d: %s", i+1, string(msgBytes))

			err = conn.WriteMessage(websocket.TextMessage, msgBytes)
			require.NoError(t, err)

			// Small delay between messages
			time.Sleep(100 * time.Millisecond)
		}

		// Read messages for 15 seconds
		done := make(chan struct{})
		go func() {
			time.Sleep(15 * time.Second)
			close(done)
		}()

	readLoop2:
		for {
			select {
			case <-done:
				break readLoop2
			default:
				conn.SetReadDeadline(time.Now().Add(2 * time.Second))
				msgType, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						t.Logf("WebSocket error: %v", err)
						break readLoop2
					}
					// Timeout or other errors, continue
					continue
				}

				if msgType != websocket.TextMessage {
					continue
				}

				var msg map[string]interface{}
				if err := json.Unmarshal(message, &msg); err != nil {
					continue
				}

				// Check if it's an update message
				if msgType, ok := msg["type"].(string); ok && msgType == "update" {
					if payload, ok := msg["payload"].(map[string]interface{}); ok {
						if symbol, ok := payload["symbol"].(string); ok {
							mu.Lock()
							switch symbol {
							case "solusdt":
								if !receivedSOL {
									receivedSOL = true
									t.Logf("✓ Received SOL update")
								}
							case "ethusdt":
								if !receivedETH {
									receivedETH = true
									t.Logf("✓ Received ETH update")
								}
							case "btcusdt":
								if !receivedBTC {
									receivedBTC = true
									t.Logf("✓ Received BTC update")
								}
							}
							mu.Unlock()
						}
					}
				}
			}
		}
		mu.Lock()
		gotSOL := receivedSOL
		gotETH := receivedETH
		gotBTC := receivedBTC
		mu.Unlock()

		t.Logf("\nResults: SOL=%v, ETH=%v, BTC=%v", gotSOL, gotETH, gotBTC)

		// Assert: Only the LAST subscription (BTC) should receive updates
		// This proves that multiple separate subscription messages don't work either
		require.False(t, gotSOL, "SOL should NOT receive updates (only last subscription works)")
		require.False(t, gotETH, "ETH should NOT receive updates (only last subscription works)")
		require.True(t, gotBTC, "BTC should receive updates (it's the last subscription)")
	})

	t.Run("SingleMessageSingleSubscriptionWithMultipleSymbolsInFilter", func(t *testing.T) {
		t.Log("=== Test: Send ONE message with ONE subscription but filters containing multiple symbols ===")

		var receivedSOL, receivedETH, receivedBTC bool
		var mu sync.Mutex

		// Connect to WebSocket
		conn, _, err := websocket.DefaultDialer.Dial("wss://ws-live-data.polymarket.com", nil)
		require.NoError(t, err)
		defer conn.Close()

		t.Log("✓ Connected")

		// Send ONE subscription message with filters containing MULTIPLE symbols
		subscribeMsg := map[string]interface{}{
			"action": "subscribe",
			"subscriptions": []map[string]string{
				{
					"topic":   "crypto_prices",
					"type":    "update",
					"filters": `{"symbols":["solusdt","ethusdt","btcusdt"]}`, // Note: "symbols" array instead of "symbol" string
				},
			},
		}

		msgBytes, _ := json.Marshal(subscribeMsg)
		t.Logf("Sending: %s", string(msgBytes))

		err = conn.WriteMessage(websocket.TextMessage, msgBytes)
		require.NoError(t, err)

		t.Log("Waiting for responses...")

		// Read messages for 15 seconds
		done := make(chan struct{})
		go func() {
			time.Sleep(15 * time.Second)
			close(done)
		}()

	readLoop3:
		for {
			select {
			case <-done:
				break readLoop3
			default:
				conn.SetReadDeadline(time.Now().Add(2 * time.Second))
				msgType, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						t.Logf("WebSocket error: %v", err)
						break readLoop3
					}
					continue
				}

				if msgType != websocket.TextMessage {
					continue
				}

				var msg map[string]interface{}
				if err := json.Unmarshal(message, &msg); err != nil {
					continue
				}

				// Check if it's an update message
				if mType, ok := msg["type"].(string); ok && mType == "update" {
					if payload, ok := msg["payload"].(map[string]interface{}); ok {
						if symbol, ok := payload["symbol"].(string); ok {
							mu.Lock()
							switch symbol {
							case "solusdt":
								if !receivedSOL {
									receivedSOL = true
									t.Logf("✓ Received SOL update")
								}
							case "ethusdt":
								if !receivedETH {
									receivedETH = true
									t.Logf("✓ Received ETH update")
								}
							case "btcusdt":
								if !receivedBTC {
									receivedBTC = true
									t.Logf("✓ Received BTC update")
								}
							}
							mu.Unlock()
						}
					}
				}
			}
		}

		mu.Lock()
		gotSOL := receivedSOL
		gotETH := receivedETH
		gotBTC := receivedBTC
		mu.Unlock()

		t.Logf("\nResults: SOL=%v, ETH=%v, BTC=%v", gotSOL, gotETH, gotBTC)

		// Assert: The server should reject the connection or send no data
		// Using symbols array is not supported by the API
		require.False(t, gotSOL, "SOL should NOT receive updates (symbols array not supported)")
		require.False(t, gotETH, "ETH should NOT receive updates (symbols array not supported)")
		require.False(t, gotBTC, "BTC should NOT receive updates (symbols array not supported)")
	})
}
