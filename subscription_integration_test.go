package polymarketrealtime

import (
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSubscriptionLifecycle tests the full subscription lifecycle: subscribe, receive messages, unsubscribe
func TestSubscriptionLifecycle(t *testing.T) {
	var receivedSubMessages []subscriptionMessage
	var sentMessages []string
	var mu sync.Mutex

	// Create a mock WebSocket server
	server, wsURL := mockWebSocketServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Handle messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			mu.Lock()
			var subMsg subscriptionMessage
			if err := json.Unmarshal(message, &subMsg); err == nil {
				receivedSubMessages = append(receivedSubMessages, subMsg)

				// Send acknowledgment
				if subMsg.Action == "subscribe" {
					ack := `{"status": "subscribed", "topic": "` + string(subMsg.Subscriptions[0].Topic) + `"}`
					conn.WriteMessage(websocket.TextMessage, []byte(ack))
					sentMessages = append(sentMessages, ack)

					// Send a test message for the subscribed topic
					testMsg := `{"topic": "` + string(subMsg.Subscriptions[0].Topic) + `", "type": "` + string(subMsg.Subscriptions[0].Type) + `", "payload": {"data": "test"}}`
					conn.WriteMessage(websocket.TextMessage, []byte(testMsg))
					sentMessages = append(sentMessages, testMsg)
				} else if subMsg.Action == "unsubscribe" {
					ack := `{"status": "unsubscribed", "topic": "` + string(subMsg.Subscriptions[0].Topic) + `"}`
					conn.WriteMessage(websocket.TextMessage, []byte(ack))
					sentMessages = append(sentMessages, ack)
				}
			}
			mu.Unlock()
		}
	})
	defer server.Close()

	var clientReceivedMessages []string
	var clientMu sync.Mutex

	// Create client
	client := New(
		WithHost(wsURL),
		WithOnNewMessage(func(data []byte) {
			clientMu.Lock()
			clientReceivedMessages = append(clientReceivedMessages, string(data))
			clientMu.Unlock()
		}),
	)

	// Connect
	err := client.Connect()
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Subscribe to activity
	subscriptions := []Subscription{
		{
			Topic: TopicActivity,
			Type:  MessageTypeAll,
		},
	}
	err = client.Subscribe(subscriptions)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// Verify subscription was sent
	mu.Lock()
	assert.GreaterOrEqual(t, len(receivedSubMessages), 1)
	assert.Equal(t, "subscribe", receivedSubMessages[0].Action)
	assert.Equal(t, TopicActivity, receivedSubMessages[0].Subscriptions[0].Topic)
	mu.Unlock()

	// Verify client received messages
	clientMu.Lock()
	assert.GreaterOrEqual(t, len(clientReceivedMessages), 2) // Ack + test message
	clientMu.Unlock()

	// Unsubscribe
	err = client.Unsubscribe(subscriptions)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// Verify unsubscribe was sent
	mu.Lock()
	assert.GreaterOrEqual(t, len(receivedSubMessages), 2)
	assert.Equal(t, "unsubscribe", receivedSubMessages[1].Action)
	assert.Equal(t, TopicActivity, receivedSubMessages[1].Subscriptions[0].Topic)
	mu.Unlock()

	// Clean up
	err = client.Disconnect()
	assert.NoError(t, err)
}

// TestMultipleSubscriptionsSameTopic tests subscribing to the same topic multiple times
func TestMultipleSubscriptionsSameTopic(t *testing.T) {
	var receivedSubMessages []subscriptionMessage
	var mu sync.Mutex

	// Create a mock WebSocket server
	server, wsURL := mockWebSocketServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Handle messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			mu.Lock()
			var subMsg subscriptionMessage
			if err := json.Unmarshal(message, &subMsg); err == nil {
				receivedSubMessages = append(receivedSubMessages, subMsg)

				// Send acknowledgment
				ack := `{"status": "ok"}`
				conn.WriteMessage(websocket.TextMessage, []byte(ack))
			}
			mu.Unlock()
		}
	})
	defer server.Close()

	// Create client
	client := New(WithHost(wsURL))

	// Connect
	err := client.Connect()
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Subscribe to the same topic multiple times with different filters
	subscription1 := []Subscription{
		{
			Topic:   TopicActivity,
			Type:    MessageTypeTrades,
			Filters: `{"market": "market1"}`,
		},
	}

	subscription2 := []Subscription{
		{
			Topic:   TopicActivity,
			Type:    MessageTypeTrades,
			Filters: `{"market": "market2"}`,
		},
	}

	subscription3 := []Subscription{
		{
			Topic:   TopicActivity,
			Type:    MessageTypeOrdersMatched,
			Filters: `{"market": "market1"}`,
		},
	}

	// Subscribe three times
	err = client.Subscribe(subscription1)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	err = client.Subscribe(subscription2)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	err = client.Subscribe(subscription3)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Verify all subscriptions were sent
	mu.Lock()
	assert.Equal(t, 3, len(receivedSubMessages))

	// Verify each subscription
	assert.Equal(t, "subscribe", receivedSubMessages[0].Action)
	assert.Equal(t, TopicActivity, receivedSubMessages[0].Subscriptions[0].Topic)
	assert.Equal(t, MessageTypeTrades, receivedSubMessages[0].Subscriptions[0].Type)
	assert.Equal(t, `{"market": "market1"}`, receivedSubMessages[0].Subscriptions[0].Filters)

	assert.Equal(t, "subscribe", receivedSubMessages[1].Action)
	assert.Equal(t, TopicActivity, receivedSubMessages[1].Subscriptions[0].Topic)
	assert.Equal(t, MessageTypeTrades, receivedSubMessages[1].Subscriptions[0].Type)
	assert.Equal(t, `{"market": "market2"}`, receivedSubMessages[1].Subscriptions[0].Filters)

	assert.Equal(t, "subscribe", receivedSubMessages[2].Action)
	assert.Equal(t, TopicActivity, receivedSubMessages[2].Subscriptions[0].Topic)
	assert.Equal(t, MessageTypeOrdersMatched, receivedSubMessages[2].Subscriptions[0].Type)
	assert.Equal(t, `{"market": "market1"}`, receivedSubMessages[2].Subscriptions[0].Filters)
	mu.Unlock()

	// Clean up
	err = client.Disconnect()
	assert.NoError(t, err)
}

// TestMultipleSubscriptionsDifferentTopics tests subscribing to different topics
func TestMultipleSubscriptionsDifferentTopics(t *testing.T) {
	var receivedSubMessages []subscriptionMessage
	var mu sync.Mutex

	// Create a mock WebSocket server
	server, wsURL := mockWebSocketServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Handle messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			mu.Lock()
			var subMsg subscriptionMessage
			if err := json.Unmarshal(message, &subMsg); err == nil {
				receivedSubMessages = append(receivedSubMessages, subMsg)

				// Send acknowledgment
				ack := `{"status": "ok"}`
				conn.WriteMessage(websocket.TextMessage, []byte(ack))
			}
			mu.Unlock()
		}
	})
	defer server.Close()

	// Create client
	client := New(WithHost(wsURL))

	// Connect
	err := client.Connect()
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Subscribe to different topics
	subscriptions := []Subscription{
		{
			Topic: TopicActivity,
			Type:  MessageTypeAll,
		},
		{
			Topic:   TopicComments,
			Type:    MessageTypeCommentCreated,
			Filters: `{"event_id": 100}`,
		},
		{
			Topic:   TopicCryptoPrices,
			Type:    MessageTypeUpdate,
			Filters: `{"symbol": "btcusdt"}`,
		},
	}

	// Subscribe all at once
	err = client.Subscribe(subscriptions)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// Verify subscription was sent with all topics
	mu.Lock()
	assert.Equal(t, 1, len(receivedSubMessages))
	assert.Equal(t, "subscribe", receivedSubMessages[0].Action)
	assert.Equal(t, 3, len(receivedSubMessages[0].Subscriptions))

	// Verify each subscription
	assert.Equal(t, TopicActivity, receivedSubMessages[0].Subscriptions[0].Topic)
	assert.Equal(t, MessageTypeAll, receivedSubMessages[0].Subscriptions[0].Type)

	assert.Equal(t, TopicComments, receivedSubMessages[0].Subscriptions[1].Topic)
	assert.Equal(t, MessageTypeCommentCreated, receivedSubMessages[0].Subscriptions[1].Type)
	assert.Equal(t, `{"event_id": 100}`, receivedSubMessages[0].Subscriptions[1].Filters)

	assert.Equal(t, TopicCryptoPrices, receivedSubMessages[0].Subscriptions[2].Topic)
	assert.Equal(t, MessageTypeUpdate, receivedSubMessages[0].Subscriptions[2].Type)
	assert.Equal(t, `{"symbol": "btcusdt"}`, receivedSubMessages[0].Subscriptions[2].Filters)
	mu.Unlock()

	// Clean up
	err = client.Disconnect()
	assert.NoError(t, err)
}

// TestUnsubscribeSpecificSubscription tests unsubscribing from a specific subscription
func TestUnsubscribeSpecificSubscription(t *testing.T) {
	var receivedSubMessages []subscriptionMessage
	var mu sync.Mutex

	// Create a mock WebSocket server
	server, wsURL := mockWebSocketServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Handle messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			mu.Lock()
			var subMsg subscriptionMessage
			if err := json.Unmarshal(message, &subMsg); err == nil {
				receivedSubMessages = append(receivedSubMessages, subMsg)

				// Send acknowledgment
				ack := `{"status": "ok"}`
				conn.WriteMessage(websocket.TextMessage, []byte(ack))
			}
			mu.Unlock()
		}
	})
	defer server.Close()

	// Create client
	client := New(WithHost(wsURL))

	// Connect
	err := client.Connect()
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Subscribe to multiple topics
	subscription1 := Subscription{
		Topic: TopicActivity,
		Type:  MessageTypeTrades,
	}

	subscription2 := Subscription{
		Topic:   TopicComments,
		Type:    MessageTypeCommentCreated,
		Filters: `{"event_id": 100}`,
	}

	err = client.Subscribe([]Subscription{subscription1, subscription2})
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Unsubscribe from only one topic
	err = client.Unsubscribe([]Subscription{subscription1})
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Verify both subscribe and unsubscribe were sent
	mu.Lock()
	assert.Equal(t, 2, len(receivedSubMessages))

	// First message should be subscribe with 2 subscriptions
	assert.Equal(t, "subscribe", receivedSubMessages[0].Action)
	assert.Equal(t, 2, len(receivedSubMessages[0].Subscriptions))

	// Second message should be unsubscribe with 1 subscription
	assert.Equal(t, "unsubscribe", receivedSubMessages[1].Action)
	assert.Equal(t, 1, len(receivedSubMessages[1].Subscriptions))
	assert.Equal(t, TopicActivity, receivedSubMessages[1].Subscriptions[0].Topic)
	assert.Equal(t, MessageTypeTrades, receivedSubMessages[1].Subscriptions[0].Type)
	mu.Unlock()

	// Clean up
	err = client.Disconnect()
	assert.NoError(t, err)
}

// TestDuplicateSubscriptions tests that duplicate subscriptions are sent to the server
func TestDuplicateSubscriptions(t *testing.T) {
	var receivedSubMessages []subscriptionMessage
	var mu sync.Mutex

	// Create a mock WebSocket server
	server, wsURL := mockWebSocketServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Handle messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			mu.Lock()
			var subMsg subscriptionMessage
			if err := json.Unmarshal(message, &subMsg); err == nil {
				receivedSubMessages = append(receivedSubMessages, subMsg)

				// Send acknowledgment
				ack := `{"status": "ok"}`
				conn.WriteMessage(websocket.TextMessage, []byte(ack))
			}
			mu.Unlock()
		}
	})
	defer server.Close()

	// Create client
	client := New(WithHost(wsURL))

	// Connect
	err := client.Connect()
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Subscribe to the exact same subscription three times
	subscription := []Subscription{
		{
			Topic:   TopicActivity,
			Type:    MessageTypeTrades,
			Filters: `{"market": "test"}`,
		},
	}

	// First subscription
	err = client.Subscribe(subscription)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Second subscription (duplicate)
	err = client.Subscribe(subscription)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Third subscription (duplicate)
	err = client.Subscribe(subscription)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Verify all three subscriptions were sent (no deduplication by default)
	mu.Lock()
	assert.Equal(t, 3, len(receivedSubMessages))

	for i := 0; i < 3; i++ {
		assert.Equal(t, "subscribe", receivedSubMessages[i].Action)
		assert.Equal(t, 1, len(receivedSubMessages[i].Subscriptions))
		assert.Equal(t, TopicActivity, receivedSubMessages[i].Subscriptions[0].Topic)
		assert.Equal(t, MessageTypeTrades, receivedSubMessages[i].Subscriptions[0].Type)
		assert.Equal(t, `{"market": "test"}`, receivedSubMessages[i].Subscriptions[0].Filters)
	}
	mu.Unlock()

	// Clean up
	err = client.Disconnect()
	assert.NoError(t, err)
}

// TestSubscriptionPersistenceAcrossReconnect tests that subscriptions are restored after reconnection
func TestSubscriptionPersistenceAcrossReconnect(t *testing.T) {
	var receivedSubMessages []subscriptionMessage
	var connectionCount int
	var mu sync.Mutex

	// Create a mock WebSocket server that closes connection after first subscription
	server, wsURL := mockWebSocketServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		mu.Lock()
		connectionCount++
		currentConnection := connectionCount
		mu.Unlock()

		// Handle messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			mu.Lock()
			var subMsg subscriptionMessage
			if err := json.Unmarshal(message, &subMsg); err == nil {
				receivedSubMessages = append(receivedSubMessages, subMsg)

				// Send acknowledgment
				ack := `{"status": "ok"}`
				conn.WriteMessage(websocket.TextMessage, []byte(ack))

				// Close connection after first subscription on first connection
				if currentConnection == 1 && subMsg.Action == "subscribe" {
					mu.Unlock()
					time.Sleep(50 * time.Millisecond)
					return
				}
			}
			mu.Unlock()
		}
	})
	defer server.Close()

	// Create client with auto-reconnect
	client := New(
		WithHost(wsURL),
		WithAutoReconnect(true),
		WithReconnectBackoff(10*time.Millisecond, 100*time.Millisecond),
	)

	// Connect
	err := client.Connect()
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Subscribe
	subscriptions := []Subscription{
		{
			Topic: TopicActivity,
			Type:  MessageTypeAll,
		},
	}

	err = client.Subscribe(subscriptions)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// Wait for reconnection (server closed the connection)
	time.Sleep(500 * time.Millisecond)

	// Verify subscriptions were sent twice (original + after reconnect)
	mu.Lock()
	assert.GreaterOrEqual(t, len(receivedSubMessages), 2)
	assert.GreaterOrEqual(t, connectionCount, 2)

	// Both should be subscribe actions
	assert.Equal(t, "subscribe", receivedSubMessages[0].Action)
	assert.Equal(t, "subscribe", receivedSubMessages[1].Action)
	mu.Unlock()

	// Clean up
	err = client.Disconnect()
	assert.NoError(t, err)
}

// TestConcurrentSubscriptions tests concurrent subscription operations
func TestConcurrentSubscriptions(t *testing.T) {
	var receivedSubMessages []subscriptionMessage
	var mu sync.Mutex

	// Create a mock WebSocket server
	server, wsURL := mockWebSocketServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Handle messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			mu.Lock()
			var subMsg subscriptionMessage
			if err := json.Unmarshal(message, &subMsg); err == nil {
				receivedSubMessages = append(receivedSubMessages, subMsg)

				// Send acknowledgment
				ack := `{"status": "ok"}`
				conn.WriteMessage(websocket.TextMessage, []byte(ack))
			}
			mu.Unlock()
		}
	})
	defer server.Close()

	// Create client
	client := New(WithHost(wsURL))

	// Connect
	err := client.Connect()
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Perform concurrent subscriptions
	var wg sync.WaitGroup
	subscriptionCount := 10

	for i := 0; i < subscriptionCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			subscription := []Subscription{
				{
					Topic:   TopicActivity,
					Type:    MessageTypeTrades,
					Filters: `{"market": "market` + string(rune('0'+index)) + `"}`,
				},
			}

			client.Subscribe(subscription)
		}(i)
	}

	wg.Wait()
	time.Sleep(200 * time.Millisecond)

	// Verify all subscriptions were sent
	mu.Lock()
	assert.Equal(t, subscriptionCount, len(receivedSubMessages))

	// All should be subscribe actions
	for _, msg := range receivedSubMessages {
		assert.Equal(t, "subscribe", msg.Action)
		assert.Equal(t, TopicActivity, msg.Subscriptions[0].Topic)
		assert.Equal(t, MessageTypeTrades, msg.Subscriptions[0].Type)
	}
	mu.Unlock()

	// Clean up
	err = client.Disconnect()
	assert.NoError(t, err)
}
