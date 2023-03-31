# How to launch a new testnet

1. Download and install `cosmos-init`: https://github.com/danil-lashin/cosmos-init
2. Edit `testnet-config.yml`
3. Run `cosmos-init testnet-config.yml`
4. Upload folders from `./chain_data` to servers.
5. Build `mineplex-chaind` and `mineplex-explorer-api` binaries on server
6. Run nodes (multiple instances):
   - `mineplex-chaind start --home ./chain_data/seed`
   - `mineplex-chaind start --home ./chain_data/validator1`
   - `mineplex-chaind start --home ./chain_data/validator...`
7. Create postgres database `mineplex`
8. Run Explorer API Service:
    - `DB_USER=postgres DB_PASSWORD=password DB_NAME=mineplex DB_HOST=127.0.0.1 DB_PORT=5432 mineplex-explorer-api`

*Note: testnet validators are configured to be able to launch on single server.*

By manipulating `testnet-config.yml` you can add/edit:
1. Accounts
2. Validators
3. All possible `genesis.json` params
4. All `config.toml`, `app.toml`, `client.toml` fields of validators and on seed node