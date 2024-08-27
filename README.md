## Avail DA Module

### Objective
The goal of the Avail DA Module is to enable Cosmos chains to leverage the Avail Network for data availability without requiring the chain to function as a rollup via the Rollkit framework. This approach provides a seamless integration for existing Layer 1 (L1) Cosmos chains to benefit from data availability, similar to integrating any other Cosmos SDK module.

### Key Elements

#### Block Interval
The `Block Interval` defines the frequency at which block data is posted to the Avail Network. For example, if the interval is set to `10`, data will be submitted at block heights `11`, `21`, `31`, and so on. At each of these intervals, the block data from the previous interval will be posted. 

For example:
- At height `11`, blocks from `1` to `10` are posted.
- At height `21`, blocks from `11` to `20` are posted.

#### Relayer
The `Relayer` acts as the transport layer, responsible for handling requests from the `prepareBlocker` and facilitating transactions between the Cosmos chain and the Avail DA network. It performs key functions such as submitting block data to Avail and updating block status on the Cosmos chain. Every validator in the network is required to run the relayer process.

#### Proven Height
The `Proven Height` represents the latest block height of the Cosmos chain for which data has been successfully posted to Avail and verified by the network.

## Architecture

1. **Block Interval Trigger**:
   - At each block interval, a request is sent from `PrepareProposal` abci method to the relayer, specifying the range of block heights to be posted to the Avail DA network. This request should be made by the block proposer only.

2. **MsgSubmitBlobRequest Transaction**:
   - The relayer initiates a `MsgSubmitBlobRequest` transaction on the Cosmos chain, marking the block data for the specified range as pending:
     ``` 
     status[range] = pending
     ```
   - The relayer tracks the transaction to ensure its successful completion.

3. **Data Submission to Avail DA**:
   - Once the `MsgSubmitBlobRequest` transaction is confirmed, the relayer fetches the block data for the specified range and submits it to the Avail DA layer.

4. **MsgUpdateBlobStatusRequest** Transaction**:
   - After confirming that the data is available on Avail, the relayer submits a `MsgUpdateBlobStatusRequest` transaction on the Cosmos chain, updating the block status to pre-verification:
     ``` 
     status[range] = pre_verification
     ```

5. **Validator Confirmation**:
   - Within a preconfigured block limit, all validators are required to verify the data's availability on the Avail network using their Avail light clients and cast their votes.

            we could use voteExtension to cast the votes
        
6. **Consensus and Proven Height Update**:
   - If the number of votes exceeds the consensus threshold, the status of the block range is updated to success, and the `Proven Height` is advanced:
     ``` 
     status[range] = success
     
     // Update the proven height
     if range.from == provenHeight + 1 {
         provenHeight = range.to
     }
     ```

7. **Failure Handling**:
   - In case of any failures or expiration of the verification window, the data will be reposted following the same procedure.

---
