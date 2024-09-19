# Cada (Cosmos Avail DA module)

CADA is a module designed to connect Cosmos sovereign chains with the Avail network, making it easier for any Cosmos chain or rollapp to use Avail as their Data Availability (DA) layer. With CADA, developers can improve the scalability and security of their decentralized applications within the Cosmos ecosystem. It enables better data handling and availability, allowing Cosmos-based chains to tap into the strengths of Avail and build a more connected and resilient blockchain network.

For example:
Let blobInterval = 10,

- At height `11`, blocks from `1` to `10` are posted.
- At height `21`, blocks from `11` to `20` are posted.

Refer to the module specification available [here](./specs/README.md) for more detailed information.

Note: Use the latest maintained [Go](https://go.dev/dl/) version to work with this module.

Ensure that the Avail light client URL is correctly configured for the module to function as expected. For instructions on running Avail locally, refer to [this documentation](https://github.com/rollkit/avail-da?tab=readme-ov-file#avail-da).
