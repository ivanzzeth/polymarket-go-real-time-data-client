package polymarketrealtime

import "encoding/json"

// ActivityFilter represents the filter configuration for activity subscriptions
type ActivityFilter struct {
	// Market is the market condition ID to filter by (optional)
	Market string `json:"market,omitempty"`
	// AssetID is the asset ID to filter by (optional)
	AssetID string `json:"asset_id,omitempty"`
}

// ToJSON converts the filter to JSON string
func (f *ActivityFilter) ToJSON() string {
	if f == nil {
		return ""
	}
	data, _ := json.Marshal(f)
	return string(data)
}

// CommentFilter represents the filter configuration for comment subscriptions
type CommentFilter struct {
	// EventID is the event ID to filter comments by (optional)
	EventID int64 `json:"event_id,omitempty"`
	// SeriesID is the series ID to filter comments by (optional)
	SeriesID int64 `json:"series_id,omitempty"`
}

// ToJSON converts the filter to JSON string
func (f *CommentFilter) ToJSON() string {
	if f == nil {
		return ""
	}
	data, _ := json.Marshal(f)
	return string(data)
}

// CryptoPriceFilter represents the filter configuration for crypto price subscriptions
type CryptoPriceFilter struct {
	// Symbol is the crypto symbol to subscribe to (e.g., "btcusdt", "ethusdt")
	Symbol string `json:"symbol"`
}

// ToJSON converts the filter to JSON string
func (f *CryptoPriceFilter) ToJSON() string {
	if f == nil {
		return ""
	}
	data, _ := json.Marshal(f)
	return string(data)
}

// EquityPriceFilter represents the filter configuration for equity price subscriptions
type EquityPriceFilter struct {
	// Symbol is the equity symbol to subscribe to (e.g., "AAPL", "GOOGL")
	Symbol string `json:"symbol"`
}

// ToJSON converts the filter to JSON string
func (f *EquityPriceFilter) ToJSON() string {
	if f == nil {
		return ""
	}
	data, _ := json.Marshal(f)
	return string(data)
}

// CLOBMarketFilter represents the filter configuration for CLOB market subscriptions
type CLOBMarketFilter struct {
	// TokenIDs is the list of token IDs to subscribe to
	TokenIDs []string `json:"token_ids"`
}

// ToJSON converts the filter to JSON string
func (f *CLOBMarketFilter) ToJSON() string {
	if f == nil || len(f.TokenIDs) == 0 {
		return ""
	}
	// The API expects a JSON array format like ["100","200"]
	data, _ := json.Marshal(f.TokenIDs)
	return string(data)
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
