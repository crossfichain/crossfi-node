# Mineplex EVM upgrade

Mineplex Chain is preparing for an upgrade to add support for EVM smart contracts. Assuming the proposal passes the chain will stop at given upgrade height.

The logs should look something like:

```bash
E[2019-11-05|12:44:18.913] UPGRADE "<plan-name>" NEEDED at height: <desired-upgrade-height>:       module=main
```

**Only after** you see this log you can shut down your node, replace binary with the new one and start it again. If you do it before upgrade height, you will permanently damage your node’s data.

You can find actual node’s binary release here: https://github.com/mineplexio-org/mineplex-2-node/releases

Notes:

1. If you are a validator and you miss upgrade block, when you need to upgrade and sync your node to latest height and then send `unjail` transaction to participate in consensus again.
2. Remember, that if you are a validator and you will occasionally run 2+ nodes with your private key (config/priv_validator_key.json), then your validator will be `tombstoned` (banned forever) and punished by 5% of its stake.
3. If you are willing to automate upgrade process you might want to setup `cosmovisor` daemon: [https://docs.cosmos.network/v0.46/run-node/cosmovisor.html](https://docs.cosmos.network/v0.46/run-node/cosmovisor.html)