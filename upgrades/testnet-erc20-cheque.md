### Node Upgrade Instructions For CrossFi Testnet

As we approach block 3,520,000, we will be performing a significant upgrade named **erc20-cheque-testnet**.
To ensure a smooth transition, please follow the detailed instructions below for upgrading your node
and using Cosmovisor to automate the upgrade process.

Assuming the proposal passes the chain will stop at given upgrade height.

The logs should look something like:

```bash
UPGRADE "erc20-cheque-testnet" NEEDED at height: 3520000
```

**Only after** you see this log you can shut down your node, replace binary with the new one and start it again. If you
do it before upgrade height, you will permanently damage your node’s data.

You can find actual node’s binary release
here: https://github.com/crossfichain/crossfi-node/releases/tag/v0.3.0-prebuild9

### Notes

1. If you are a validator, and you miss upgrade block, when you need to upgrade and sync your node to latest height and
   then send `unjail` transaction to participate in consensus again.
2. Remember, that if you are a validator, and you will occasionally run 2+ nodes with your private key (
   config/priv_validator_key.json), then your validator will be `tombstoned` (banned forever) and punished by 5% of its
   stake.
3. If you are willing to automate upgrade process you might want to set up `cosmovisor`
   daemon: [https://docs.cosmos.network/v0.46/run-node/cosmovisor.html](https://docs.cosmos.network/v0.46/run-node/cosmovisor.html)

#### Setting up Cosmovisor for automated upgrade

Please refer to [Cosmovisor documentation](https://docs.cosmos.network/main/build/tooling/cosmovisor) to setup automatic
upgrade for your node.