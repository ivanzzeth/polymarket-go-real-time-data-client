package polymarketrealtime

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ActivityFilter represents the filter configuration for activity subscriptions
type ActivityFilter struct {
	// Market is the market condition ID to filter by (optional)
	Market string `json:"market,omitempty"`
	// AssetID is the asset ID to filter by (optional)
	AssetID string `json:"asset_id,omitempty"`
}

// ToJSON converts the filter to JSON string
// Returns error if filter is invalid (both Market and AssetID are empty)
func (f *ActivityFilter) ToJSON() (string, error) {
	if f == nil {
		return "", nil
	}
	if f.Market == "" && f.AssetID == "" {
		return "", fmt.Errorf("ActivityFilter: at least one of Market or AssetID must be specified")
	}
	data, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("ActivityFilter: failed to marshal to JSON: %w", err)
	}
	return string(data), nil
}

// CommentFilter represents the filter configuration for comment subscriptions
type CommentFilter struct {
	// EventID is the event ID to filter comments by (optional)
	EventID int64 `json:"event_id,omitempty"`
	// SeriesID is the series ID to filter comments by (optional)
	SeriesID int64 `json:"series_id,omitempty"`
}

// ToJSON converts the filter to JSON string
// Returns error if filter is invalid (both EventID and SeriesID are zero)
func (f *CommentFilter) ToJSON() (string, error) {
	if f == nil {
		return "", nil
	}
	if f.EventID == 0 && f.SeriesID == 0 {
		return "", fmt.Errorf("CommentFilter: at least one of EventID or SeriesID must be specified")
	}
	data, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("CommentFilter: failed to marshal to JSON: %w", err)
	}
	return string(data), nil
}

// CryptoPriceFilter represents the filter configuration for crypto price subscriptions
type CryptoPriceFilter struct {
	// Symbol is the crypto symbol to subscribe to (e.g., "btcusdt", "ethusdt")
	Symbol string `json:"symbol"`
}

// ToJSON converts the filter to JSON string
// Returns error if filter is invalid (symbol is empty)
func (f *CryptoPriceFilter) ToJSON() (string, error) {
	if f == nil {
		return "", fmt.Errorf("CryptoPriceFilter: filter cannot be nil")
	}
	if f.Symbol == "" {
		return "", fmt.Errorf("CryptoPriceFilter: symbol must be specified")
	}
	// Validate symbol format (should be lowercase)
	if strings.ToLower(f.Symbol) != f.Symbol {
		return "", fmt.Errorf("CryptoPriceFilter: symbol must be lowercase (e.g., 'btcusdt' not 'BTCUSDT')")
	}
	data, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("CryptoPriceFilter: failed to marshal to JSON: %w", err)
	}
	return string(data), nil
}

// EquityPriceFilter represents the filter configuration for equity price subscriptions
type EquityPriceFilter struct {
	// Symbol is the equity symbol to subscribe to (e.g., "AAPL", "GOOGL")
	Symbol string `json:"symbol"`
}

// ToJSON converts the filter to JSON string
// Returns error if filter is invalid (symbol is empty)
func (f *EquityPriceFilter) ToJSON() (string, error) {
	if f == nil {
		return "", fmt.Errorf("EquityPriceFilter: filter cannot be nil")
	}
	if f.Symbol == "" {
		return "", fmt.Errorf("EquityPriceFilter: symbol must be specified")
	}
	// Validate symbol format (should be uppercase)
	if strings.ToUpper(f.Symbol) != f.Symbol {
		return "", fmt.Errorf("EquityPriceFilter: symbol must be uppercase (e.g., 'AAPL' not 'aapl')")
	}
	data, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("EquityPriceFilter: failed to marshal to JSON: %w", err)
	}
	return string(data), nil
}

// CLOBMarketFilter represents the filter configuration for CLOB market subscriptions
type CLOBMarketFilter struct {
	// TokenIDs is the list of token IDs to subscribe to
	TokenIDs []string `json:"token_ids"`
}

// ToJSON converts the filter to JSON string
// Returns error if filter is invalid (token IDs are empty)
func (f *CLOBMarketFilter) ToJSON() (string, error) {
	if f == nil {
		return "", fmt.Errorf("CLOBMarketFilter: filter cannot be nil")
	}
	if len(f.TokenIDs) == 0 {
		return "", fmt.Errorf("CLOBMarketFilter: at least one token ID must be specified")
	}
	// Validate token IDs are not empty strings
	for i, tokenID := range f.TokenIDs {
		if tokenID == "" {
			return "", fmt.Errorf("CLOBMarketFilter: token ID at index %d is empty", i)
		}
	}
	// The API expects a JSON array format like ["100","200"]
	data, err := json.Marshal(f.TokenIDs)
	if err != nil {
		return "", fmt.Errorf("CLOBMarketFilter: failed to marshal to JSON: %w", err)
	}
	return string(data), nil
}

// NewActivityFilter creates a new activity filter
func NewActivityFilter() *ActivityFilter {
	return &ActivityFilter{}
}

// WithMarket sets the market filter
func (f *ActivityFilter) WithMarket(market string) *ActivityFilter {
	f.Market = market
	return f
}

// WithAssetID sets the asset ID filter
func (f *ActivityFilter) WithAssetID(assetID string) *ActivityFilter {
	f.AssetID = assetID
	return f
}

// NewCommentFilter creates a new comment filter
func NewCommentFilter() *CommentFilter {
	return &CommentFilter{}
}

// WithEventID sets the event ID filter
func (f *CommentFilter) WithEventID(eventID int64) *CommentFilter {
	f.EventID = eventID
	return f
}

// WithSeriesID sets the series ID filter
func (f *CommentFilter) WithSeriesID(seriesID int64) *CommentFilter {
	f.SeriesID = seriesID
	return f
}

// NewCryptoPriceFilter creates a new crypto price filter with the given symbol
func NewCryptoPriceFilter(symbol string) *CryptoPriceFilter {
	return &CryptoPriceFilter{Symbol: symbol}
}

// NewEquityPriceFilter creates a new equity price filter with the given symbol
func NewEquityPriceFilter(symbol string) *EquityPriceFilter {
	return &EquityPriceFilter{Symbol: symbol}
}

// NewCLOBMarketFilter creates a new CLOB market filter with the given token IDs
func NewCLOBMarketFilter(tokenIDs ...string) *CLOBMarketFilter {
	return &CLOBMarketFilter{TokenIDs: tokenIDs}
}

// AddTokenID adds a token ID to the filter
func (f *CLOBMarketFilter) AddTokenID(tokenID string) *CLOBMarketFilter {
	f.TokenIDs = append(f.TokenIDs, tokenID)
	return f
}
