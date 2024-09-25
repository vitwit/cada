<!--
order: 5
-->

# ProofOfBlobProposalHandler: PrepareProposal Method

This documentation provides an overview of the `PrepareProposal` method within the `ProofOfBlobProposalHandler`. This method is critical for preparing a block proposal by aggregating and injecting vote information into the proposal transactions.

## Method Overview

The `PrepareProposal` method performs the following key steps:

### 1. Proposer Address Initialization

The method starts by setting the `proposerAddress` in the keeper with the address provided in the `RequestPrepareProposal`. This address represents the proposer of the current block. Since the `PrepareProposal` ABCI method is exclusively executed by the block proposer, this address can later be used for posting block data to `Avail`.

```go
h.keeper.proposerAddress = req.ProposerAddress
```


### 2. Vote Aggregation

The method then aggregates votes by calling the `aggregateVotes` function, which takes in the current context and the `LocalLastCommit` from the request. This function collects votes from the last commit, which are essential for the consensus process for da verification.

```go
votes, err := h.aggregateVotes(ctx, req.LocalLastCommit)
if err != nil {
    fmt.Println("error while aggregating votes", err)
    return nil, err
}
```



### 3. Injection of Aggregated Votes

The method creates a new structure, `StakeWeightedVotes`, to hold the aggregated votes and the extended commit information (`ExtendedCommitInfo`).

```go
injectedVoteExtTx := StakeWeightedVotes{
    Votes:              votes,
    ExtendedCommitInfo: req.LocalLastCommit,
}
bz, err := json.Marshal(injectedVoteExtTx)
if err != nil {
    fmt.Println("failed to encode injected vote extension tx", "err", err)
}
```

The serialized vote information (`injectedVoteExtTx`) is appended to the list of proposal transactions, which can be later processed in `PreBlocker` abci method.

```go
proposalTxs = append(proposalTxs, bz)

return &abci.ResponsePrepareProposal{
    Txs: proposalTxs,
}, nil
```