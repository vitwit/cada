# Configuration 

The avail-sdk configuration is located in config/app.toml and is used for connecting to the avail network.

Below is the default configuration for integrating with Avail. You can customize these settings by replacing them with your own light client url and Avail seed.

[avail]

    # avail light client node url for posting data
    light-client-url = "http://127.0.0.1:8000"

    # Avail chain id
    chain-id = "avail-1"

    # Overrides the expected chain's app-id, test-only
    override-app-id = "1"

    # Overrides the expected chain's publish-to-avali light client block interval, test-only
    override-pub-interval = 5

    # Query avail for new block proofs this often
    proof-query-interval = "12s"

    # Seed for avail
    seed = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"

    # RPC of cosmos node to get the block data
    cosmos-node-rpc = "http://127.0.0.1:26657"