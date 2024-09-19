# CADA (Cosmos Avail DA module)

CADA is a module designed to connect Cosmos sovereign chains with the Avail network, making it easier for any Cosmos chain or rollapp to use Avail as their Data Availability (DA) layer. With CADA, developers can improve the scalability and security of their decentralized applications within the Cosmos ecosystem. It enables better data handling and availability, allowing Cosmos-based chains to tap into the strengths of Avail and build a more connected and resilient blockchain network.

# How It Works
For example:
Let blobInterval = 10,

- At height `11`, blocks from `1` to `10` are posted.
- At height `21`, blocks from `11` to `20` are posted.

For more detailed information, refer to the `CADA` module specification [here](./specs/README.md).

# Integration Guide

To integrate the CADA module into your application, follow the steps outlined in the [integration guide](./integration_docs/README.md)

Note: Ensure that the Avail light client URL is correctly configured for the module to function as expected. For instructions on setup Avail locally, please refer to [this documentation](https://github.com/rollkit/avail-da?tab=readme-ov-file#avail-da).
