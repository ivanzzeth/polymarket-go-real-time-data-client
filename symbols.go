package polymarketrealtime

// Crypto symbols for price subscriptions
const (
	// Major cryptocurrencies
	CryptoSymbolBTCUSDT = "btcusdt" // Bitcoin
	CryptoSymbolETHUSDT = "ethusdt" // Ethereum
	CryptoSymbolSOLUSDT = "solusdt" // Solana
	CryptoSymbolBNBUSDT = "bnbusdt" // Binance Coin
	CryptoSymbolXRPUSDT = "xrpusdt" // Ripple
	CryptoSymbolADAUSDT = "adausdt" // Cardano
	CryptoSymbolDOGEUSDT = "dogeusdt" // Dogecoin
	CryptoSymbolMATICUSDT = "maticusdt" // Polygon
	CryptoSymbolDOTUSDT = "dotusdt" // Polkadot
	CryptoSymbolAVAXUSDT = "avaxusdt" // Avalanche
	CryptoSymbolSHIBUSDT = "shibusdt" // Shiba Inu
	CryptoSymbolLTCUSDT = "ltcusdt" // Litecoin
	CryptoSymbolLINKUSDT = "linkusdt" // Chainlink
	CryptoSymbolUNIUSDT = "uniusdt" // Uniswap
	CryptoSymbolATOMUSDT = "atomusdt" // Cosmos
	CryptoSymbolETCUSDT = "etcusdt" // Ethereum Classic
	CryptoSymbolXLMUSDT = "xlmusdt" // Stellar
	CryptoSymbolNEARUSDT = "nearusdt" // NEAR Protocol
	CryptoSymbolAPTUSDT = "aptusdt" // Aptos
	CryptoSymbolARBUSDT = "arbusdt" // Arbitrum
	CryptoSymbolOPUSDT = "opusdt" // Optimism
)

// Equity symbols for price subscriptions
const (
	// Technology stocks
	EquitySymbolAAPL = "AAPL" // Apple Inc.
	EquitySymbolMSFT = "MSFT" // Microsoft Corporation
	EquitySymbolGOOGL = "GOOGL" // Alphabet Inc. (Google)
	EquitySymbolAMZN = "AMZN" // Amazon.com Inc.
	EquitySymbolNVDA = "NVDA" // NVIDIA Corporation
	EquitySymbolTSLA = "TSLA" // Tesla Inc.
	EquitySymbolMETA = "META" // Meta Platforms Inc. (Facebook)
	EquitySymbolNFLX = "NFLX" // Netflix Inc.
	EquitySymbolAMD = "AMD" // Advanced Micro Devices
	EquitySymbolINTC = "INTC" // Intel Corporation
	EquitySymbolCSCO = "CSCO" // Cisco Systems Inc.
	EquitySymbolORCL = "ORCL" // Oracle Corporation
	EquitySymbolADBE = "ADBE" // Adobe Inc.
	EquitySymbolCRM = "CRM" // Salesforce Inc.
	EquitySymbolPYPL = "PYPL" // PayPal Holdings Inc.

	// Financial stocks
	EquitySymbolJPM = "JPM" // JPMorgan Chase & Co.
	EquitySymbolBAC = "BAC" // Bank of America Corp.
	EquitySymbolWFC = "WFC" // Wells Fargo & Company
	EquitySymbolGS = "GS" // Goldman Sachs Group Inc.
	EquitySymbolMS = "MS" // Morgan Stanley
	EquitySymbolV = "V" // Visa Inc.
	EquitySymbolMA = "MA" // Mastercard Inc.

	// Healthcare stocks
	EquitySymbolJNJ = "JNJ" // Johnson & Johnson
	EquitySymbolUNH = "UNH" // UnitedHealth Group Inc.
	EquitySymbolPFE = "PFE" // Pfizer Inc.
	EquitySymbolABBV = "ABBV" // AbbVie Inc.
	EquitySymbolMRK = "MRK" // Merck & Co. Inc.

	// Consumer stocks
	EquitySymbolWMT = "WMT" // Walmart Inc.
	EquitySymbolHD = "HD" // Home Depot Inc.
	EquitySymbolDIS = "DIS" // Walt Disney Company
	EquitySymbolNKE = "NKE" // Nike Inc.
	EquitySymbolKO = "KO" // Coca-Cola Company
	EquitySymbolPEP = "PEP" // PepsiCo Inc.
	EquitySymbolMCD = "MCD" // McDonald's Corporation

	// Energy stocks
	EquitySymbolXOM = "XOM" // Exxon Mobil Corporation
	EquitySymbolCVX = "CVX" // Chevron Corporation

	// Industrial stocks
	EquitySymbolBA = "BA" // Boeing Company
	EquitySymbolCAT = "CAT" // Caterpillar Inc.
	EquitySymbolGE = "GE" // General Electric Company
)

// Convenience functions for creating filters with predefined symbols

// NewBTCPriceFilter creates a filter for Bitcoin price
func NewBTCPriceFilter() *CryptoPriceFilter {
	return NewCryptoPriceFilter(CryptoSymbolBTCUSDT)
}

// NewETHPriceFilter creates a filter for Ethereum price
func NewETHPriceFilter() *CryptoPriceFilter {
	return NewCryptoPriceFilter(CryptoSymbolETHUSDT)
}

// NewSOLPriceFilter creates a filter for Solana price
func NewSOLPriceFilter() *CryptoPriceFilter {
	return NewCryptoPriceFilter(CryptoSymbolSOLUSDT)
}

// NewAppleStockFilter creates a filter for Apple stock price
func NewAppleStockFilter() *EquityPriceFilter {
	return NewEquityPriceFilter(EquitySymbolAAPL)
}

// NewTeslaStockFilter creates a filter for Tesla stock price
func NewTeslaStockFilter() *EquityPriceFilter {
	return NewEquityPriceFilter(EquitySymbolTSLA)
}

// NewNvidiaStockFilter creates a filter for NVIDIA stock price
func NewNvidiaStockFilter() *EquityPriceFilter {
	return NewEquityPriceFilter(EquitySymbolNVDA)
}
