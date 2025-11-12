package polymarketrealtime

import (
	"encoding/json"
	"time"
)

// "connection_id": "",
// "timestamp":1762928533586,
//    "topic":"activity",
//    "type":"orders_matched"

type SubscriptionMessage struct {
	ConnectionID string          `json:"connection_id"`
	Timestamp    int64           `json:"timestamp"`
	Time         time.Time       `json:"-"` // Parsed from timestamp
	Topic        Topic           `json:"topic"`
	Type         MessageType     `json:"type"`
	Payload      json.RawMessage `json:"payload"`
}

// MarshalJSON implements the json.Marshaler interface for Message
func (m SubscriptionMessage) MarshalJSON() ([]byte, error) {
	type Alias SubscriptionMessage
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for Message
func (m *SubscriptionMessage) UnmarshalJSON(data []byte) error {
	type Alias SubscriptionMessage
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse timestamp to time.Time
	if m.Timestamp != 0 {
		// Assuming timestamp is in milliseconds
		m.Time = time.Unix(0, m.Timestamp*int64(time.Millisecond))
	}

	return nil
}
