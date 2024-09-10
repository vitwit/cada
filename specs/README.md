
# `x/cada`


## Abstract

This document specifies the cada module of the Cosmos SDK.

CADA is a module designed to connect Cosmos sovereign chains with the Avail network, making it easier for any Cosmos chain or rollapp to use Avail as their Data Availability (DA) layer. With CADA, developers can improve the scalability and security of their decentralized applications within the Cosmos ecosystem. It enables better data handling and availability, allowing Cosmos-based chains to tap into the strengths of Avail and build a more connected and resilient blockchain network.



### MaxBlocksLimitForBlob

* The `maxBlocksLimitForBlob` defines the maximum number blobs to be posted to Avail DA at once (in on Avail Transaction).

### Blob Interval
* The `Blob Interval` defines the frequency at which block data is posted to the Avail Network. 
* For example, if the interval is set to `5`, data will be submitted at block heights `6`, `11`, `16`, and so on. 
* At each of these intervals, the block data from the proven height to min(current height, proven height + maxBlocksLimitForBlob) will be posted. 

Example:
For Blob Interval = 5 and Maximum Blocks Limit for Blob = 10 :-

- At height `6` and provenHeight = `0`, blocks from `1` to `5` are posted.

- At height `11` and provenHeight still `0`, blocks from `1` to `10` are posted.

### Relayer
* The `Relayer` acts as the transport layer, responsible for handling requests from the `prepareBlocker` and facilitating transactions between the Cosmos chain and the Avail DA network. 
* It performs key functions such as submitting block data to Avail and updating block status on the Cosmos chain. Every validator in the network is required to run the relayer process.
* Relayer should initialized with a chain account so that the validator can use this account to sign `MsgUpdateStatusBlob` transaction.

### Voting Interval
* The `Voting Interval` is the period before validators verify whether data is truly included in Avail and confirm it with the network using vote extensions.


## State

The module keeps state of the following primary objects:

1. **Blocks Height Range**: Tracks the start and end of the current blocks range being posted to Avail.
2. **Blob Submission Status**: Indicates the current status of the blob submission process (`READY`, `PENDING`, `IN_VOTING`, `FAILURE`).
3. **Voting End Height**: The block height at which the voting for the current blocks should conclude.
4. **Avail Height**: The Avail block height at which the data is made available.
5. **Proven Height**: The latest block height of the Cosmos chain for which data has been successfully posted to Avail and verified by the network.

The module uses the following keys to manage the aforementioned state:

* **Height Range Start Key**: `0x07` - Stores the start of the range of current blocks being posted to Avail.
* **Height Range End Key**: `0x08` - Stores the end of the range of current blocks being posted to Avail.
* **Blob Status Key**: `0x06` - Stores the status of the blob submission process.
* **Voting End Height Key**: `0x09` - Stores the block height at which voting should end for the current blocks.
* **Avail Height Key**: `0x0A` - Stores the Avail block height where the data is made available.



## Architecture


   - At each block interval, a request is sent from the `PreBlocker` ABCI method to the Keeper, specifying the range of block heights that are ready to be posted to the `Avail` DA network.
   - The range of block heights should be from `provenHeight + 1` to `min(provenHeight + MaxBlocksLimitForBlob, CurrentBlockHeight)`.

   - If the status of the previous blocks is either `READY` or `FAILURE`, the status can be updated to `PENDING`.
     
     ``` 
     range = [fromBlock, toBlock] // (fromBlock < toBlock < CurrentBlock)
     status = PENDING
     ```

    - The `Proposer` of the block will make a request to the `Relayer` to post the blocks data by passing the range of blocks to be posted.

    - The `Relayer` fetches the blocks data from the local provider, converts the blocks data to bytes, and posts that data to `Avail`.

    - Once the success of data availability is confirmed, the `Relayer` broadcasts the `Avail height` at which the blob data is made available using the `MsgUpdateBlobStatus` transaction.

    - The status, Avail height, and voting deadline will be updated in the state.

    ```
    status = IN_VOTING
    availHeight = tx.availHeight
    votingEndBlock = currentBlock + votingInterval
    ```

    - At block height `VotingEndBlock - 1`, all the validators verify if the specified blocks data is truly made available at the specified Avail height. They cast their vote (YES or NO) using `vote extensions`.

    - At block height `VotingEndBlock`, all the votes from `vote_extensions` will be collected and aggregated. If the collective `voting power is > 66%`, the status will be updated

    ```
    status = READY // success and ready for next blocks
    provenHeight = Range End

    ```
    - In case of failure at any stage, the whole flow will be repeated.


---
