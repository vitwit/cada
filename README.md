# Cada (Cosmos Avail DA module)

CADA is a module designed to connect Cosmos sovereign chains with the Avail network, making it easier for any Cosmos chain or rollapp to use Avail as their Data Availability (DA) layer. With CADA, developers can improve the scalability and security of their decentralized applications within the Cosmos ecosystem. It enables better data handling and availability, allowing Cosmos-based chains to tap into the strengths of Avail and build a more connected and resilient blockchain network.

For example:
Let blobInterval = 10,
- At height `11`, blocks from `1` to `10` are posted.
- At height `21`, blocks from `11` to `20` are posted.

#### Relayer
The `Relayer` acts as the transport layer, responsible for handling requests from the `prepareBlocker` and facilitating transactions between the Cosmos chain and the Avail DA network. It performs key functions such as submitting block data to Avail and updating block status on the Cosmos chain. Every validator in the network is required to run the relayer process.

#### Proven Height
The `Proven Height` signifies the most recent block height of the Cosmos chain where data has been successfully transmitted to Avail and validated by the network.

## Architecture

![Screenshot from 2024-08-27 11-35-01](https://github.com/user-attachments/assets/1a8657f6-4c1b-418a-8295-05c039baa6d0)


1. **Block Interval Trigger**:
   - At each block interval, a request is sent from `PrepareProposal` abci method to the relayer, specifying the range of block heights to be posted to the Avail DA network. This request should be made by the block proposer only.

2. **MsgSubmitBlobRequest Transaction**:
   - The relayer submits a `MsgSubmitBlobRequest` transaction on the Cosmos chain, signaling that the block data for the specified range is pending:
     ``` 
     status[range] = pending
     ```
   - The relayer monitors the transaction to confirm its successful inclusion and processing on the chain.

3. **Data Submission to Avail DA**:
   - Once the `MsgSubmitBlobRequest` transaction is confirmed, the relayer fetches the block data for the specified range and submits it to the Avail DA layer.

4. **MsgUpdateBlobStatusRequest Transaction**:
   - After confirming that the data is available on Avail, the relayer submits a `MsgUpdateBlobStatusRequest` transaction on the Cosmos chain, updating the block status to pre-verification:
     ``` 
     status[range] = IN_VOTING
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
For detailed instructions on how to integrate the module with a spawn generated application, please refer to the [integration guide](./docs/spawn.md).
