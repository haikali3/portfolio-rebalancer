Test Coverage (README says this is evaluation criteria)

  - No handler/integration tests for /portfolio or /rebalance
  - No storage layer tests
  - No Kafka producer/consumer tests
  - No edge case tests (empty allocations, missing user, invalid inputs)

  Fault Tolerance (README explicitly calls this out)

  - No retry logic on consumer message processing failures
  - No dead letter queue for failed messages
  - No Elasticsearch volume mount (data lost on restart)

  Missing Features

  - No GET endpoints to retrieve portfolios or transactions
  - No health check endpoint
  - No allocation percentage validation (e.g., should sum to 100%)
  - Portfolio isn't updated after rebalance — only transactions are saved

  Code Quality / Production Readiness

  - No structured logging (just log.Printf)
  - No input validation middleware
  - No graceful shutdown handling
  - README TODO says "Write a README" — the current one is the assignment prompt, not a project README
