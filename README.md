# Substate 
Substate database is used as the Off-The-Chain testing module in the applications for recording or replaying transactions. The Replayer can execute any transaction in complete isolation by loading the substate of the transaction and executing the transaction.

Database contains a minimal subset of the World-State Trie to faithfully replay transactions in isolation. The subset contains all the entries represented as a flat key-value store (and is not stored as a slow Merkle Patricia Trie) for executing a transaction.

## Data Structure

There are 5 data structures stored in a substate DB:
1. `SubstateAccount`: account information (nonce, balance, code, storage)
2. `WorldState`: mapping of account address and `SubstateAccount`
3. `SubstateEnv`: information from block headers (block gasLimit, number, timestamp, hashes)
4. `SubstateMessage`: message for transaction execution
5. `SubstateResult`: result of transaction execution

5 values are required to replay transactions and validate results:
1. `PreState`: world-state that is read during transaction execution
2. `Env`: block information required for transaction execution
3. `Message`: array with exactly 1 transaction
4. `PostState`: alloc that is generated by transaction execution
5. `Result`: execution result and receipt array with exactly 1 receipt

The first 2 bytes of a key in a substate DB represent different data types as follows:
1. `1s`: Substate, a key is `"1s"+N+T` with transaction index `T` at block `N`.
`T` and `N` are encoded in a big-endian 64-bit binary.
2. `1c`: EVM bytecode, a key is `"1c"+codeHash` where `codeHash` is Keccak256 hash of the bytecode.

# Ethereum Substate Recorder/Replayer
Ethereum substate recorder/replayer based on the paper:

**Yeonsoo Kim, Seongho Jeong, Kamil Jezek, Bernd Burgstaller, and Bernhard Scholz**: _An Off-The-Chain Execution Environment for Scalable Testing and Profiling of Smart Contracts_,  USENIX ATC'21

You can find all executables including `geth` and our `substate-cli` in `build/bin/` directory.