```bash

mineplex-chaind init mynode
cd ~/.mineplex-chain/config/
wget https://raw.githubusercontent.com/mineplexio/MinePlex-2/master/networks/testnet/genesis.json

mineplex-chaind start --p2p.persistent_peers="fd12d31114a0bf44872f1584bd8fb02bdbda6503@172.104.249.147:26656"

```