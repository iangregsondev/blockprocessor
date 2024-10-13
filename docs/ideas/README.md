# Ideas and considerations

## Retry situations

The current provider that I am using for the blockchain data as a rate limiter, so I do get some HTTP 429 errors.
This is circumvented by retrying the request after a delay. The delay is not increasing, it is a fixed delay.
This is configurable in the configuration files

As we are using kafka topics if the processing of messages fails at any point and because we are using consumer groups, the message will be reprocessed.

## Loosing transactions

If the service is down for an hour, the messages will be stored in the kafka topic.
When the service is restarted, the messages will be reprocessed.

There is a retention on the kafka topics, so if the service is down for more than the retention period, the messages will be lost.

## Block reorganisation

This would be an issue on Bitcoin, in bitcoin it is consensus that a block that has > 6 confirmations is considered final.
We could only process blocks that have > 6 confirmations, this would ensure that the block is final and thus the transactions in the block are final.

I think this would be the best way to handle this situation.

Otherwise, the alternative option is a lot more complex.

We could introduce another topic, the blocks will still get added to the current topic in realtime but an additional service could
check blocks that have greater than 6 confirmations and move these another topic, this new topic is where the transaction processor would read from.

There are many ways to handle this, but I think the two options I have listed are the best.





