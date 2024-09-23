# Configuration 

The cada module configuration is located in config/app.toml and is used for connecting to the avail and cosmos networks to submit blocks data.

Below is the default configuration for integrating with Avail. You can customize these settings by modifying them to suit your specific requirements.

[avail]

    # Avail light client node url for posting data
    light-client-url = "http://127.0.0.1:8000"

    # Avail chain id
    chain-id = "avail-1"

    # Overrides the expected  avail app-id, test-only
    override-app-id = "1"

    # Seed for avail
    seed = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"

    # RPC of cosmos node to get the block data
    cosmos-node-rpc = "http://127.0.0.1:26657"

    # Maximum number of blocks over which blobs can be processed
    max-blob-blocks = 10

    # The frequency at which block data is posted to the Avail Network
    publish-blob-interval = 5

    # It is the period before validators verify whether data is truly included in
    # Avail and confirm it with the network using vote extension
    vote-interval = 5

    # It is the keyname of the cosmos validator account to sign the transactions
    validator-key = "alice"