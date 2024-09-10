<!--
order: 2
-->

# State

The module keeps state of the following primary objects:

## Blocks Range

Tracks the start and end of the current blocks range being posted to `Avail`.

### Blocks Range Start Height

Stores the start height of the range of current blocks being posted to `Avail`.

It is stored in the state as follows:

- PrevHeightKey `0x07` -> Start Height (Uint64)

### Blocks Range End Height

Stores the end height of the range of current blocks being posted to `Avail`.

It is stored in the state as follows:

- NextHeightKey `0x08` -> End Height (Uint64)

## Blocks Submission Status

Indicates the status of the current blocks submission (`READY`, `PENDING`, `IN_VOTING`, `FAILURE`).

** PENDING ** : Blocks data submission has been initiated and is awaiting confirmation
** IN_VOTING ** : Blocks data has been posted to `Avail` and is now pending validators' verification
** FAILURE ** : Blocks data submission or verification has failed and needs to be resubmitted
** READY ** : blocks data submission is successful; the next set of blocks is ready to be posted

It is stored in the state as follows:

- BlobStatusKey `0x06` : status (uint32)

## Voting End Height

The block height at which the voting for the current blocks should conclude.

It is stored in the state as follows:

- VotingEndHeightKey `0x09` : voting end block height (uint64)

## Avail Height

The Avail block height at which the data is made available.

It is stored in the state as follows:

AvailHeightKey `0x0A` : avail block height (uint64)

## Proven Height

The latest block height of the Cosmos chain for which data has been successfully posted to Avail and verified by the network.

It is stored in the state as follows:

ProvenHeightKey `0x02` : proven block height (uint64)
