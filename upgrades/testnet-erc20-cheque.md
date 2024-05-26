### Node Upgrade Instructions For CrossFi Testnet

As we approach block ___, we will be performing a significant upgrade named **erc20-cheque-testnet**.
To ensure a smooth transition, please follow the detailed instructions below for upgrading your node
and using Cosmovisor to automate the upgrade process.

Assuming the proposal passes the chain will stop at given upgrade height.

The logs should look something like:

```bash
E[2019-11-05|12:44:18.913] UPGRADE "<plan-name>" NEEDED at height: <desired-upgrade-height>:       module=main
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

*Prerequisites*

1. **Cosmovisor**: Ensure you have Cosmovisor installed. If not, you can install it following the instructions in
   the [Cosmovisor documentation](https://docs.cosmos.network/main/build/tooling/cosmovisor).
2. Here we assume that your node is being run by `systemd`. If it is not - please consider your own setup.

1. **Stop the Node**
   ```bash
   sudo systemctl stop crossfid
   ```

2. **Install the New Version**
    - Download the new binary corresponding to the **erc20-cheque-testnet** upgrade.
   ```bash
   https://github.com/crossfichain/crossfi-node/releases/download/v0.3.0-prebuild9/crossfi-node_0.3.0-prebuild9_linux_amd64.tar.gz
   tar -xvf crossfi-node_0.3.0-prebuild9_linux_amd64.tar.gz
   ```

3. **Prepare Cosmovisor Directories**
    - Set up the Cosmovisor directory structure:
   ```bash
   mkdir -p $HOME/.crossfid/cosmovisor/genesis/bin
   mkdir -p $HOME/.crossfid/cosmovisor/upgrades/erc20-cheque-testnet/bin
   ```

    - Copy the current and new binaries into the respective directories:
   ```bash
   cp $(which crossfid) $HOME/.crossfid/cosmovisor/genesis/bin/
   cp ./bin/crossfid $HOME/.crossfid/cosmovisor/upgrades/erc20-cheque-testnet/bin/
   ```

4. **Update the Systemd Service File**
    - Edit the systemd service file to use Cosmovisor:
   ```bash
   sudo nano /etc/systemd/system/crossfid.service
   ```

    - Modify the `ExecStart` line to use Cosmovisor:
   ```plaintext
   [Unit]
   Description=Crossfid Daemon
   After=network-online.target

   [Service]
   User=<your-username>
   ExecStart=/usr/local/bin/cosmovisor run start
   Restart=always
   RestartSec=3
   LimitNOFILE=4096
   Environment="DAEMON_NAME=crossfid"
   Environment="DAEMON_HOME=$HOME/.crossfid"
   Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=true"
   Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
   Environment="UNSAFE_SKIP_BACKUP=true"

   [Install]
   WantedBy=multi-user.target
   ```

    - Save and close the file. Then reload systemd:
   ```bash
   sudo systemctl daemon-reload
   ```

5. **Start the Node with Cosmovisor**
   ```bash
   sudo systemctl start crossfid
   ```

6. **Monitor the Upgrade**
    - Watch the logs to ensure the node is running properly and the upgrade occurs at block 100:
   ```bash
   journalctl -fu crossfid
   ```

#### Post-Upgrade Verification

1. **Check Node Status**
    - Verify the node is running the new version:
   ```bash
   $HOME/.crossfid/cosmovisor/current/bin/crossfid version
   ```

2. **Sync Status**
    - Ensure the node is fully synced:
   ```bash
   $HOME/.crossfid/cosmovisor/current/bin/crossfid status
   ```

By following these instructions, you can ensure that your node transitions smoothly to the **erc20-cheque-testnet**
upgrade using Cosmovisor. If you encounter any issues during the process, please reach out to the support team or
consult the community forums for assistance.

Thank you for your attention and cooperation.

---

**Note**: Replace placeholders like `<your-username>` with appropriate values specific to your environment.