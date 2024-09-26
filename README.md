# CADA (Cosmos Avail DA module)

CADA is a module designed to connect Cosmos sovereign chains with the Avail network, making it easier for any Cosmos chain or rollapp to use Avail as their Data Availability (DA) layer. With CADA, developers can improve the scalability and security of their decentralized applications within the Cosmos ecosystem. It enables better data handling and availability, allowing Cosmos-based chains to tap into the strengths of Avail and build a more connected and resilient blockchain network.

# Integration Guide

To integrate the CADA module into your application, follow the steps outlined in the [integration guide](./integration_docs/README.md)

Note: Ensure that the Avail light client URL is correctly configured for the module to function as expected. For instructions on setup Avail locally, please refer to [this documentation](https://github.com/rollkit/avail-da?tab=readme-ov-file#avail-da).

# How It Works

`Cosmos Chain`: Initiates the process by running the PreBlocker ABCI method.

`Request to Relayer`: Sends block range information to the relayer.

`Relayer`: Fetches the block data from the Cosmos Provider and posts it to the Avail Light Client.

`Avail Light Client`: Confirms whether the data is available.

`If Yes`: Broadcast the Avail height and status.

`If No`: Retry data submission.

`Validators`: Vote to confirm the data availability, updating the status to "Ready" or "Failure" based on results.

There are main components in the workflow:

## 1. Cada
The core functionality of the **Cada** module is integrated with and operates on the Cosmos blockchain.

In the Cada module:
- At each block interval, the `PreBlocker` ABCI method sends a request to the `Relayer`, specifying the range of block heights that are ready to be posted to the **Avail** Data Availability (DA) network.
![Data Submission](https://github.com/user-attachments/assets/fc4d23cc-f6bd-4210-8407-47a57adcc290)

- The chain is responsible for aggregating vote extensions from all validators and verifying whether the data has been made available on Avail.
- Since verification requires communicating with the light client, an asynchronous voting mechanism is needed. **Vote extensions** enable this asynchronous voting mechanism for verification purposes.

![Vote Extension](https://github.com/user-attachments/assets/ea5b10ab-fb64-4ed0-8761-44675a852a01)

## 2. Relayer
The **Relayer** facilitates communication between the Cosmos Chain, the Avail light client, and the Cosmos Provider.

- **Data Submission**: The relayer is responsible for fetching block data from the cosmos provider and posting it to the Avail light client via an HTTP request.
- Based on the response from the light client, the relayer submits a transaction informing the validators of the data availability status and the specific Avail block height where the data is included, so that validators can verify it.
  
- **Data Verification**: During verification, the relayer communicates with the Avail light client to confirm whether the data is truly available at the specified height.



## 3. Avail Light Node
The **Avail Light Client** allows interaction with the Avail DA network without requiring a full node, and without having to trust remote peers. It leverages **Data Availability Sampling (DAS)**, which the light client performs on every newly created block.

- The chain communicates with the Avail light client via the relayer during the data submission and data availability verification processes.

Find more details about the Avail Light Client [here](https://docs.availproject.org/docs/operate-a-node/run-a-light-client/Overview).

## 4. Cosmos Provider
The **Cosmos Provider** is responsible for fetching block data via RPC so that the data can be posted to Avail for availability checks.


# Architecture


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

