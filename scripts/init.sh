#!/bin/bash

# Clean existing data
rm -rf ~/.dhivesd

# Initialize chain
dhivesd init testnode --chain-id dhives_testnet-1 --home ~/.dhivesd

# Create keys
dhivesd keys add validator1 --keyring-backend test --home ~/.dhivesd
dhivesd keys add validator2 --keyring-backend test --home ~/.dhivesd
dhivesd keys add validator3 --keyring-backend test --home ~/.dhivesd

# Add genesis accounts
dhivesd add-genesis-account $(dhivesd keys show validator1 -a --keyring-backend test --home ~/.dhivesd) 4000000000000000000000000dhives --home ~/.dhivesd
dhivesd add-genesis-account $(dhivesd keys show validator2 -a --keyring-backend test --home ~/.dhivesd) 3000000000000000000000000dhives --home ~/.dhivesd
dhivesd add-genesis-account $(dhivesd keys show validator3 -a --keyring-backend test --home ~/.dhivesd) 3000000000000000000000000dhives --home ~/.dhivesd

# Create validator transaction
dhivesd gentx validator1 1000000000000000000000000dhives   --chain-id dhives_testnet-1   --moniker="validator1"   --commission-max-change-rate=0.01   --commission-max-rate=0.20   --commission-rate=0.10   --min-self-delegation=1   --keyring-backend test   --home ~/.dhivesd

# Collect genesis transactions
dhivesd collect-gentxs --home ~/.dhivesd

# Validate genesis
dhivesd validate-genesis --home ~/.dhivesd

# Update configurations
if [ -f ~/.dhivesd/config/app.toml ]; then
    sed -i '' 's/minimum-gas-prices = "0stake"/minimum-gas-prices = "0.1dhives"/' ~/.dhivesd/config/app.toml
    sed -i '' 's/enable = false/enable = true/' ~/.dhivesd/config/app.toml
fi

echo "Chain initialized! Start with: dhivesd start --home ~/.dhivesd"
