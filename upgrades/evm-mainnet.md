# Instructions for Validators to Migrate to CrossFi EVM Mainnet

This guide will help you migrate your validator to the new **CrossFi EVM** chain. Follow these steps carefully to avoid
downtime and ensure a successful migration.

---

### Prerequisites:

- Ensure you have access to the old validator data, especially the `config/priv_validator_key.json`.
- Set up a new server or service.
- Download the CrossFi binaries (link provided below).

---

### Step-by-Step Migration Instructions

#### 1. **Backup Old Validator Data**

Before starting the migration, ensure you back up all essential data from your current setup. This includes:

- **`config/priv_validator_key.json` file**: This key uniquely identifies your validator. Ensure it's securely backed
  up.
- **Old chain directory**: Backup the `.mineplex-chaind` or `~/.your_old_chain_directory` that holds important data from
  the old chain.

#### 2. **Set Up a New Server or Service**

To avoid downtime during migration, set up a new server or cloud instance to run CrossFi.

- Ensure the server meets the hardware requirements for running CrossFi.

#### 3. **Setup CrossFi EVM**

```bash
wget https://github.com/crossfichain/crossfi-node/releases/download/v0.3.0/crossfi-node_0.3.0_linux_amd64.tar.gz && tar -xf crossfi-node_0.3.0_linux_amd64.tar.gz
git clone https://github.com/crossfichain/mainnet.git
```

#### 4. **Move Your `priv_validator_key.json` to CrossFi**

Your validator’s identity is tied to the `config/priv_validator_key.json`. You need to transfer this key to the new
CrossFi chain’s config directory.

#### 5. **Launch the CrossFi Validator**

Now, it's time to launch your validator on the CrossFi chain.

- Start the node with:
  ```bash
  ./bin/crossfid start --home ./mainnet
  ```

- Check the logs to ensure the node is running and syncing with the CrossFi blockchain.

#### 8. **Verify Validator Status**

Once your node is live, verify that your validator is successfully migrated.

- Check the node status:
  ```bash
  ./bin/crossfid status --home ./mainnet
  ```

- Ensure your validator is active and bonded on the CrossFi chain. You can also verify this on the CrossFi block
  explorer or by querying the chain:
  ```bash
  ./bin/crossfid --home ./mainnet query staking validators
  ```

---

### Important Notes:

- **Retain old data**: Keep your old chain’s data backed up, especially the `priv_validator_key.json` and state data.
- **Security**: Ensure that your `priv_validator_key.json` is stored securely and is never exposed.
- **Backup your setup**: Regularly back up your new CrossFi chain setup, especially the `priv_validator_key.json` and
  configuration files.

---

By following these steps, you should successfully migrate your validator to the CrossFi chain while retaining your
validator identity and data.