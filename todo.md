add volume for elasticsearch data

1. kafka integration
/rebalance handler publish msg to kafka topic
create consumer that reads from kafka
consumer runs CalculateRebalance + SaveRebalanceTransactions
handle retries on failure
start consumer in main.go

2. test
- CalculateRebalance unit test
- handler test
- edge cases: empty allocation, missing user, invalid data
- integration test with kafka and elasticsearch (if time permits)

3. write readme
- how to run
- api docs with examples
- design decisions and tradeoffs
