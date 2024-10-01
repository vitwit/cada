<!--
order: 1
-->

# Concepts

### MaxBlocksLimitForBlob

- The `maxBlocksLimitForBlob` defines the maximum number blobs to be posted to Avail DA at once (in on Avail Transaction).

### Blob Interval

- The `Blob Interval` defines the frequency at which block data is posted to the Avail Network.
- For example, if the interval is set to `5`, data will be submitted at block heights `6`, `11`, `16`, and so on.
- At each of these intervals, the block data from the proven height to min(current height, proven height + maxBlocksLimitForBlob) will be posted.

Example:
For Blob Interval = 5 and Maximum Blocks Limit for Blob = 10 :-

- At height `6` and provenHeight = `0`, blocks from `1` to `5` are posted.

- At height `11` and provenHeight still `0`, blocks from `1` to `10` are posted.

### Relayer

- The `Relayer` acts as the transport layer, responsible for handling requests from the `preBlocker` and facilitating transactions between the Cosmos chain and the Avail DA network.
- It performs key functions such as submitting block data to Avail and updating block status on the Cosmos chain. Every validator in the network is required to run the relayer process.
- Relayer should initialized with a chain account so that the validator can use this account to sign `MsgUpdateStatusBlob` transaction.

### Voting Interval

- The `Voting Interval` is the period before validators verify whether data is truly included in Avail and confirm it with the network using vote extensions.
