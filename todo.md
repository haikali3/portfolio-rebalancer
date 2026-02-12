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

Critical (must have):
1. CalculateRebalance — correct BUY/SELL transactions generated
2. CalculateRebalance — correct percentages calculated
3. HandlePortfolio — valid request returns 201
4. HandleRebalance — valid request publishes to Kafka and returns 200

Error handling:
5. HandlePortfolio — empty user_id returns 400
6. HandlePortfolio — empty allocation returns 400
7. HandlePortfolio — invalid JSON returns 400
8. HandleRebalance — empty user_id returns 400
9. HandleRebalance — empty new_allocation returns 400
10. HandleRebalance — invalid JSON returns 400

Edge cases:
11. CalculateRebalance — no change needed (same allocation) returns empty
12. CalculateRebalance — only one asset changes
13. CalculateRebalance — new asset appears in updated but not in original
14. SaveRebalanceTransactions — empty transactions slice

Not critical (nice to have):
15. HandlePortfolio — wrong HTTP method returns 405
16. HandleRebalance — wrong HTTP method returns 405
17. Consumer — unmarshal failure is handled gracefully