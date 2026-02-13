# Portfolio Rebalancer

Backend API service for managing user portfolios and calculating rebalance transactions. When a third-party provider reports that market conditions have shifted a user's allocation, the system computes the BUY/SELL transactions needed to restore the original target allocation.

## Architecture

```
Client ──POST /portfolio──▶ API Server ──▶ Elasticsearch (portfolios index)

Client ──POST /rebalance──▶ API Server ──▶ Kafka (rebalance topic)
                                                  │
                                          Consumer (with retries)
                                                  │
                                           ┌──────┴──────┐
                                           │  Calculate   │
                                           │  Rebalance   │
                                           └──────┬──────┘
                                                  │
                                    Elasticsearch (rebalance_transactions index)
                                                  │
                                         (on failure after 3 retries)
                                                  │
                                           Kafka DLQ topic
```

Rebalance requests are processed asynchronously via Kafka to handle high throughput from providers. Failed messages are retried with exponential backoff and sent to a dead letter queue after exhausting retries.

## Tech Stack

- **Go** - API server and business logic
- **Elasticsearch** - Persistence for portfolios and rebalance transactions
- **Kafka** - Message queue for async rebalance processing
- **Docker Compose** - Container orchestration

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Running

```bash
docker compose build
docker compose up
```

The service will be available at `http://localhost:8080`.

### Tooling UIs

| Tool | URL | Purpose |
|------|-----|---------|
| Kibana | http://localhost:5601 | Browse Elasticsearch data |
| Kafbat UI | http://localhost:9000 | Monitor Kafka topics and messages |

## API Reference

### POST /portfolio

Create a user portfolio with target allocation percentages.

**Request:**
```json
{
  "user_id": "1",
  "allocation": {
    "stocks": 60,
    "bonds": 30,
    "gold": 10
  }
}
```

**Response (201):**
```json
{
  "status_code": 201,
  "data": {
    "user_id": "1",
    "allocation": {
      "stocks": 60,
      "bonds": 30,
      "gold": 10
    }
  }
}
```

**Validation:**
- `user_id` is required
- `allocation` must be non-empty
- Allocation percentages must be non-negative and sum to 100

### POST /rebalance

Submit an updated allocation (from market changes) to trigger rebalancing back to the user's original target.

**Request:**
```json
{
  "user_id": "1",
  "new_allocation": {
    "stocks": 70,
    "bonds": 20,
    "gold": 10
  }
}
```

**Response (200):**
```json
{
  "status_code": 200,
  "msg": "rebalance request received and being processed"
}
```

The request is published to Kafka and processed asynchronously. The consumer fetches the user's original allocation from Elasticsearch, calculates the required transactions, and saves them.

**Example:** If the original allocation is `{stocks: 60, bonds: 30, gold: 10}` and the market has shifted it to `{stocks: 70, bonds: 20, gold: 10}`, the system generates:

| Action | Asset | Percent |
|--------|-------|---------|
| SELL | stocks | 10% |
| BUY | bonds | 10% |

## Project Structure

```
cmd/api/main.go              # Entry point, server setup
internal/
  handlers/
    portfolio.go             # HTTP handlers for /portfolio and /rebalance
    helper.go                # Error types and JSON response helpers
  services/
    rebalance.go             # Rebalance calculation logic
    rebalance_test.go        # Unit tests for rebalance calculations
  kafka/
    producer.go              # Kafka producer with DLQ support
    consumer.go              # Consumer with retry logic (3 attempts, exponential backoff)
  models/
    portfolio.go             # Portfolio, UpdatedPortfolio, RebalanceTransaction
  storage/
    elastic.go               # Elasticsearch operations (save/get portfolio, bulk save transactions)
```

## Fault Tolerance

- **Async processing** - Rebalance requests go through Kafka, decoupling the API from processing.
- **Retry logic** - Failed messages are retried up to 3 times with exponential backoff (1s, 2s, 3s).
- **Dead letter queue** - Messages that fail all retries are published to a `rebalance-dlq` topic for investigation.
- **Connection retries** - Both Elasticsearch and Kafka connections retry on startup until available.
- **Persistent storage** - Elasticsearch data is stored on a Docker volume (`es-data`).

## Running Tests

```bash
go test ./...
```
