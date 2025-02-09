#!/bin/bash

# Clean up any existing data
rm -rf ~/.dhivesd

# Initialize chain
dhivesd init testnode --chain-id dhives_testnet-1 --home ~/.dhivesd

# Verify genesis.json exists
if [ ! -f ~/.dhivesd/config/genesis.json ]; then
    echo "Error: genesis.json not found after initialization"
    exit 1
fi

# Configure keyring-backend
KEYRING="--keyring-backend test"
CHAINID="dhives_testnet-1"
HOMEDIR="--home ~/.dhivesd"

# Create keys with explicit parameters
echo "Creating keys..."
# dhivesd keys add validator1 --algo secp256k1 $KEYRING $HOMEDIR
# dhivesd keys add validator2 --algo secp256k1 $KEYRING $HOMEDIR
# dhivesd keys add validator3 --algo secp256k1 $KEYRING $HOMEDIR

# Store validator1 address
VAL1_ADDR="dhives1tx7vxgw2eg76j0uz70ts6kc3wgc59qtqqdxpuh"
VAL2_ADDR="dhives1uar9j9yj440a8gjzjglrv0ww47svrxwq95hjn7"
VAL3_ADDR="dhives1d3htlru0vtekd2yd6fm7mdng435r65rvllvuls"

# Update configurations
if [ -f ~/.dhivesd/config/app.toml ]; then
    sed -i '' 's/minimum-gas-prices = "0stake"/minimum-gas-prices = "0.1dhives"/' ~/.dhivesd/config/app.toml
    sed -i '' 's/enable = false/enable = true/' ~/.dhivesd/config/app.toml
fi

# Add genesis accounts
echo "Adding genesis accounts..."
dhivesd add-genesis-account $VAL1_ADDR 4000000000000000000000000dhives $HOMEDIR
dhivesd add-genesis-account $VAL2_ADDR 3000000000000000000000000dhives $HOMEDIR
dhivesd add-genesis-account $VAL3_ADDR 3000000000000000000000000dhives $HOMEDIR

# Create validator transaction
echo "Creating gentx..."
dhivesd gentx $HOMEDIR validator1 1000000000000000000000000dhives \
  --chain-id=$CHAINID \
  --moniker="validator1" \
  --commission-max-change-rate=0.01 \
  --commission-max-rate=0.20 \
  --commission-rate=0.10 \
  --min-self-delegation=1 \
  $KEYRING \
  $HOMEDIR

# Collect genesis transactions
echo "Collecting gentxs..."
dhivesd collect-gentxs $HOMEDIR

# Validate genesis
echo "Validating genesis..."
dhivesd validate-genesis $HOMEDIR

echo "Chain initialized! Start with: dhivesd start --home ~/.dhivesd"
