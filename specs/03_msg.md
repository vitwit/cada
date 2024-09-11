<!--
order: 3
-->

# Messages

## UpdateBlobStatus

The `MsgUpdateBlobStatus` is used to update the status of block submissions from `PENDING` to either `IN_VOTING` if the submission is successful, or `FAILURE` if it fails. The responsibility for executing this transaction lies with the individual who originally submitted the blocks data to `Avail` (the proposer of the block where the blocks data submission was initiated).

This message will fail under the following conditions:

    If the status is nil, meaning it is neither true nor false.
    If the status is true but the Avail height is not a valid number.