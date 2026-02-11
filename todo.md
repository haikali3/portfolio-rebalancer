1. POST /portfolio
- Purpose: save original allocation
- Request JSON
- Success response (201)
- Validation errors (400)

2. POST /rebalance
- Purpose: submit updated allocation for rebalance processing
- Request JSON
- Expected response (200 or 202 if async via Kafka)
- Not found/validation errors

Then do Data Flow:

- /portfolio -> Elasticsearch
- /rebalance -> Kafka -> consumer -> Elasticsearch