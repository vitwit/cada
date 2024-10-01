<!--
order: 4
-->

# Client

A user can query and interact with the `cada` module using the CLI.


## Query

The `query` commands allows users to query `cada` state.

```sh
$ simd query cada --help
```

### Query the Status

The `get-da-status` command enables users to retrieve comprehensive information, including the range of blocks currently being posted to Avail, the current status, the last proven height, the Avail height where the data is made available, and the voting block height by which voting should conclude.

```sh
$ simd query cada get-da-status
```

Output:

```yml
last_blob_voting_ends_at: "23"
proven_height: "0"
range:
  from: "1"
  to: "5"
status: IN_VOTING
```
