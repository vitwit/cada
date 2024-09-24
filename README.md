# CADA (Cosmos Avail DA module)

CADA is a module designed to connect Cosmos sovereign chains with the Avail network, making it easier for any Cosmos chain or rollapp to use Avail as their Data Availability (DA) layer. With CADA, developers can improve the scalability and security of their decentralized applications within the Cosmos ecosystem. It enables better data handling and availability, allowing Cosmos-based chains to tap into the strengths of Avail and build a more connected and resilient blockchain network.

# How It Works
For example:
Let blobInterval = 10,

- At height `11`, blocks from `1` to `10` are posted.
- At height `21`, blocks from `11` to `20` are posted.

For more detailed information, refer to the `CADA` module specification [here](./specs/README.md).

# Integration Guide

To integrate the CADA module into your application, follow the steps outlined in the [integration guide](./integration_docs/README.md)

Note: Ensure that the Avail light client URL is correctly configured for the module to function as expected. For instructions on setup Avail locally, please refer to [this documentation](https://github.com/rollkit/avail-da?tab=readme-ov-file#avail-da).


# Architecture

![blocks-data-submission](https://github.com/user-attachments/assets/4e17b98f-ca8c-4b4c-a79e-8c60f123cb2c)


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

![vote-extension](https://github.com/user-attachments/assets/c0edb8e7-20fd-468a-9109-4f31718e4467)

- At block height `VotingEndBlock - 1`, all the validators verify if the specified blocks data is truly made available at the specified Avail height. They cast their vote (YES or NO) using `vote extensions`.

- At block height `VotingEndBlock`, all the votes from `vote_extensions` will be collected and aggregated. If the collective `voting power is > 66%`, the status will be updated

    ```
    status = READY // success and ready for next blocks
    provenHeight = Range End

    ```
- In case of failure at any stage, the whole flow will be repeated.

---

