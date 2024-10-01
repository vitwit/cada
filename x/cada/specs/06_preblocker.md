<!--
order: 6
-->

# ProofOfBlobProposalHandler: PreBlocker Method

This documentation provides a detailed overview of the `PreBlocker` method within the `ProofOfBlobProposalHandler`. This method is crucial for processing vote extensions, updating statuses, and initiating block data submissions to Avail.

## Method Overview

The `PreBlocker` method is responsible for two primary tasks: processing votes from vote extensions to update the status and initiating the submission of blocks' data to Avail. Below is a step-by-step explanation of how this method works.

### 1. Process Votes from Vote Extensions and Update the Status

When the current block height matches the voting end height and the status is `IN_VOTING`, the method processes the voting results and updates the state accordingly.

- **Success Condition:** If the collective voting power exceeds 66%, the status is updated to `READY`, and the `provenHeight` is set to the end of the current block range.
- **Failure Condition:** If the voting power is 66% or less, the status is updated to `FAILURE`.

```go
if len(req.Txs) > 0 && currentHeight == int64(votingEndHeight) && blobStatus == IN_VOTING_STATE {
    var injectedVoteExtTx StakeWeightedVotes
    if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
        fmt.Println("preblocker failed to decode injected vote extension tx", "err", err)
    } else {
        from := k.GetStartHeightFromStore(ctx)
        to := k.GetEndHeightFromStore(ctx)

        pendingRangeKey := Key(from, to)
        votingPower := injectedVoteExtTx.Votes[pendingRangeKey]

        if votingPower > 66 { // Voting power is greater than 66%
            k.setBlobStatusSuccess(ctx)
        } else {
            k.SetBlobStatusFailure(ctx)
        }
    }
}
```

### 2. Initiate Block Data Availability (DA) Submission

If the current block height aligns with a voting interval and the status is either `READY` or `FAILURE`, the method updates the block range and sets the status to `PENDING` for the next round of blocks data submission.

- **Range Calculation:** The pending block range to be submitted is calculated based on the last proven height and the current block height.
- **Status Update:** The status is set to `PENDING` to mark the start of the data submission process.

```go
// The following code is executed at block heights that are multiples of the voteInterval,
// i.e., voteInterval+1, 2*voteInterval+1, 3*voteInterval+1, etc.
if !k.IsValidBlockToPostTODA(uint64(currentBlockHeight)) {
    return nil
}

provenHeight := k.GetProvenHeightFromStore(ctx)
fromHeight := provenHeight + 1                                                     // Calculate pending range of blocks to post data
endHeight := min(fromHeight+uint64(k.MaxBlocksForBlob), uint64(ctx.BlockHeight())) // Exclusive range i.e., [fromHeight, endHeight)

sdkCtx := sdk.UnwrapSDKContext(ctx)
ok := k.SetBlobStatusPending(sdkCtx, fromHeight, endHeight-1)
if !ok {
    return nil
}
```

### 3. Proposer Submits Block Data to Avail DA

If the node running this method is the proposer of the block, it takes responsibility for submitting the blocks data to Avail DA.

- **Block Data Submission:** The proposer gathers the blocks within the calculated range and posts them to Avail DA.

```go
var blocksToSubmit []int64

for i := fromHeight; i < endHeight; i++ {
    blocksToSubmit = append(blocksToSubmit, int64(i))
}

// Only the proposer should execute the following code
if bytes.Equal(req.ProposerAddress, k.proposerAddress) {
    k.relayer.PostBlocks(ctx, blocksToSubmit, k.cdc, req.ProposerAddress)
}
```