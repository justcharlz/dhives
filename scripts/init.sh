# init.sh
#!/bin/bash

# Clean existing data
rm -rf ~/.dhivesd

# Initialize chain
dhivesd init testnode --chain-id dhives_testnet-1

# Create keys
dhivesd keys add validator1 --keyring-backend test
dhivesd keys add validator2 --keyring-backend test
dhivesd keys add validator3 --keyring-backend test

# Add genesis accounts
dhivesd add-genesis-account $(dhivesd keys show validator1 -a --keyring-backend test) 4000000000000000000000000dhives
dhivesd add-genesis-account $(dhivesd keys show validator2 -a --keyring-backend test) 3000000000000000000000000dhives
dhivesd add-genesis-account $(dhivesd keys show validator3 -a --keyring-backend test) 3000000000000000000000000dhives

# Create validator transaction
dhivesd gentx validator1 1000000000000000000000000dhives \
  --chain-id dhives_testnet-1 \
  --moniker="validator1" \
  --commission-max-change-rate=0.01 \
  --commission-max-rate=0.20 \
  --commission-rate=0.10 \
  --min-self-delegation=1 \
  --keyring-backend test

# Collect genesis transactions
dhivesd collect-gentxs

# Validate genesis
dhivesd validate-genesis

# Update configurations
sed -i '' 's/minimum-gas-prices = "0stake"/minimum-gas-prices = "0.1dhives"/' ~/.dhivesd/config/app.toml
sed -i '' 's/enable = false/enable = true/' ~/.dhivesd/config/app.toml

echo "Chain initialized! Start with: dhivesd start"