# Integration Tests

This document describes the integration tests for the Polymarket Go Real-Time Data Client.

## Overview

Integration tests connect to the real Polymarket WebSocket endpoints to verify that the client works correctly with actual data streams. These tests are separate from unit tests and are tagged with the `integration` build tag.

## Running Integration Tests

### Run All Integration Tests

```bash
go test -tags=integration -v -timeout 5m
```

### Run Specific Integration Tests

#### Test Real Connection (Subscribe, Unsubscribe, Disconnect)
```bash
go test -tags=integration -v -run TestIntegrationRealConnection -timeout 2m
```

#### Test Multiple Subscriptions to Same Topic
```bash
go test -tags=integration -v -run TestIntegrationMultipleSubscriptionsSameTopic -timeout 2m
```

#### Test Multiple Different Topics
```bash
go test -tags=integration -v -run TestIntegrationMultipleTopics -timeout 2m
```

#### Test CLOB Market Subscriptions
```bash
go test -tags=integration -v -run TestIntegrationCLOBMarketSubscription -timeout 2m
```

#### Test Resubscribe After Unsubscribe
```bash
go test -tags=integration -v -run TestIntegrationResubscribeAfterUnsubscribe -timeout 2m
```

#### Test Concurrent Subscriptions
```bash
go test -tags=integration -v -run TestIntegrationConcurrentSubscriptions -timeout 2m
```

## Test Coverage

The integration tests cover the following scenarios:

### 1. Basic Connection Lifecycle
- **TestIntegrationRealConnection**
  - Connects to Polymarket WebSocket
  - Subscribes to crypto price updates (BTC)
  - Receives real-time messages
  - Unsubscribes from updates
  - Disconnects cleanly

### 2. Multiple Subscriptions to Same Topic
- **TestIntegrationMultipleSubscriptionsSameTopic**
  - Tests subscribing to the same topic (BTC prices) multiple times
  - Verifies that all subscriptions are sent to the server
  - Checks message reception behavior with duplicate subscriptions

### 3. Multiple Different Topics
- **TestIntegrationMultipleTopics**
  - Subscribes to multiple crypto symbols (BTC, ETH, SOL)
  - Uses message router to separate messages by symbol
  - Counts messages received for each symbol
  - Verifies that all subscriptions work concurrently

### 4. CLOB Market Data
- **TestIntegrationCLOBMarketSubscription**
  - Connects to CLOB WebSocket endpoint
  - Subscribes to market price changes
  - Receives and parses PriceChanges messages
  - Tests with real token IDs from active markets

### 5. Unsubscribe and Resubscribe
- **TestIntegrationResubscribeAfterUnsubscribe**
  - Phase 1: Subscribe and receive messages
  - Phase 2: Unsubscribe and verify message flow stops
  - Phase 3: Resubscribe and verify messages resume
  - Compares message counts across all phases

### 6. Concurrent Operations
- **TestIntegrationConcurrentSubscriptions**
  - Performs concurrent subscription operations
  - Subscribes to multiple cryptos simultaneously
  - Tests thread safety of the client
  - Verifies all concurrent subscriptions succeed

## Important Notes

1. **Network Required**: These tests require an active internet connection to Polymarket's WebSocket endpoints.

2. **Market Activity**: Some tests depend on market activity. If markets are inactive, you may receive fewer messages than expected.

3. **Token IDs**: CLOB market tests use real token IDs. If these markets become inactive, you may need to update the token IDs in the test code.

4. **Timeouts**: Tests include appropriate timeouts (typically 2-5 minutes) to allow for connection establishment and message reception.

5. **Rate Limiting**: Running tests too frequently may trigger rate limiting. Allow appropriate cooldown between test runs.

6. **Build Tags**: Integration tests use the `integration` build tag to prevent them from running during normal `go test` execution.

## Debugging Integration Tests

### Enable Verbose Logging

All integration tests use `WithLogger(NewLogger())` which outputs debug information. Watch the test output for:
- Connection status (✅ Connected)
- Subscription confirmations (✅ Subscribed)
- Message reception (📨 Received message)
- Disconnect status (❌ Disconnected)

### Common Issues

1. **Connection Timeout**: Check network connectivity and firewall settings
2. **No Messages Received**: Verify that the market/token is currently active
3. **Subscription Errors**: Check that the filter format is correct for the topic

## Running Specific Test Scenarios

### Test Only Real-Time Data Endpoint
```bash
go test -tags=integration -v -run "TestIntegration(RealConnection|MultipleTopics|Resubscribe|Concurrent)" -timeout 5m
```

### Test Only CLOB Endpoint
```bash
go test -tags=integration -v -run "TestIntegrationCLOB" -timeout 2m
```

## CI/CD Integration

To run integration tests in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run Integration Tests
  run: go test -tags=integration -v -timeout 5m
  env:
    GO111MODULE: on
```

Note: Consider running integration tests separately from unit tests, possibly on a schedule rather than on every commit.

## Updating Tests

When updating integration tests:

1. Verify token IDs are still valid for CLOB tests
2. Ensure timeouts are appropriate for current network conditions
3. Update assertions based on actual message formats from the API
4. Test with different market conditions (high/low activity)

## Unit Tests vs Integration Tests

- **Unit Tests** (`*_test.go` without build tags): Use mock servers, fast, run on every commit
- **Integration Tests** (`integration_test.go` with `+build integration`): Use real connections, slower, run periodically

To run only unit tests (excluding integration):
```bash
go test -v
```
