<!--
order: 7
-->

# Vote Extensions

Vote extensions are used to propagate arbitrary data across the network without needing to implement transactions that modify the state. Validators utilize vote extensions to verify data availability of a specific range of blocks and update the state of the Cada module accordingly.

This specification details the functionality and purpose of the `ExtendVoteHandler` and `VerifyVoteExtensionHandler` methods within the VoteExtHandler struct. These methods are part of a voting extension process where validators extend their votes based on the availability of data in the Avail DA network.

### ExtendVoteHandler

The `ExtendVoteHandler` method is responsible for generating a `vote extension`. It checks the availability of specific data in the Avail by interacting with an Avail light client. The vote extension is then created based on the outcome of this check.

* The method first begins by retrieving several voting-related parameters, including the start and end heights of the blocks being processed, the Avail block height, and the status of the current blob (data)

```go
from := h.Keeper.GetStartHeightFromStore(ctx)
		end := h.Keeper.GetEndHeightFromStore(ctx)

		availHeight := h.Keeper.GetAvailHeightFromStore(ctx)

		pendingRangeKey := Key(from, end)

		blobStatus := h.Keeper.GetBlobStatus(ctx)
		currentHeight := ctx.BlockHeight()
		voteEndHeight := h.Keeper.GetVotingEndHeightFromStore(ctx)
```

* The method checks if the current height is just before the end of the voting period and if the blob is in the voting state. If not, it generates a basic vote extension indicating that no data was verified

```go
if currentHeight+1 != int64(voteEndHeight) || blobStatus != IN_VOTING_STATE {
			voteExt := VoteExtension{
				Votes: Votes,
			}

			// json marshalling
			votesBytes, err := json.Marshal(voteExt)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
			}
			abciResponseVoteExt.VoteExtension = votesBytes
			return abciResponseVoteExt, nil
		}
```

* If the conditions are met, the method queries the Avail light client to determine if the relevant data is available. The result of this check is recorded in a map
```go
    ok, err := h.Keeper.relayer.IsDataAvailable(ctx, from, end, availHeight, "http://localhost:8000")
		if ok {
			h.logger.Info("submitted data to Avail verified successfully at",
				"block_height", availHeight,
			)
		}

		
		Votes[pendingRangeKey] = ok
```

* The outcome (whether the data was available or not) is marshaled into a vote extension and returned as part of the abci.ResponseExtendVote

```go
    voteExt := VoteExtension{
			Votes: Votes,
		}

		votesBytes, err := json.Marshal(voteExt)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
		}

		return &abci.ResponseExtendVote{
			VoteExtension: votesBytes,
		}, nil
```

### VerifyVoteExtensionHandler

The `VerifyVoteExtensionHandler` method is responsible for validating the format and content of the vote extensions generated by the `ExtendVoteHandler`.

This method performs a basic validation check on the received vote extension, ensuring it meets the necessary format requirements. It then returns a response indicating whether the vote extension is accepted.
