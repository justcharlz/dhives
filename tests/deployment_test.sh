#!/bin/bash

# Enable strict mode
set -euo pipefail

# Logging function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Binary configuration
BINARY="evmosd"
KEYRING="test"
KEYRING_FLAG="--keyring-backend ${KEYRING}"

# Environment variables
CHAIN_ID="evmos_9000-1"
MONIKER="localvalidator"
KEY_NAME="validator"
KEY_NAME2="user"
TOKEN_DENOM="aevmos"
NODE_HOME="$HOME/.${BINARY}"
STAKE_AMOUNT="1000000000000000000000${TOKEN_DENOM}" # 1000 aevmos with 18 decimals
MIN_GAS_PRICES="0.0001${TOKEN_DENOM}"
COMMUNITY_TOKENS="100000000000000000000000${TOKEN_DENOM}" # 100,000 aevmos

# Cleanup function
cleanup() {
    log "Cleaning up previous deployment..."
    rm -rf "${NODE_HOME}"
    pkill -f ${BINARY} || true
}

# Set up trap for cleanup on script exit
trap cleanup EXIT INT TERM

# Initialize the chain
init_chain() {
    log "Initializing chain..."
    ${BINARY} init "${MONIKER}" --chain-id "${CHAIN_ID}" --home "${NODE_HOME}"
    
    # Update stake denom to aevmos in genesis
    sed -i.bak 's/"stake"/"'${TOKEN_DENOM}'"/g' "${NODE_HOME}/config/genesis.json"
    
    # Update minimum gas prices in app.toml
    sed -i.bak "s/^minimum-gas-prices = .*/minimum-gas-prices = \"${MIN_GAS_PRICES}\"/" "${NODE_HOME}/config/app.toml"
}

# Create keys for testing
create_keys() {
    log "Creating keys..."
    ${BINARY} keys add "${KEY_NAME}" ${KEYRING_FLAG}
    ${BINARY} keys add "${KEY_NAME2}" ${KEYRING_FLAG}
}

# Set up genesis accounts
setup_genesis_accounts() {
    log "Setting up genesis accounts..."
    VALIDATOR_ADDR=$(${BINARY} keys show "${KEY_NAME}" -a ${KEYRING_FLAG})
    USER_ADDR=$(${BINARY} keys show "${KEY_NAME2}" -a ${KEYRING_FLAG})
    
    ${BINARY} add-genesis-account "${VALIDATOR_ADDR}" "${STAKE_AMOUNT}" --home "${NODE_HOME}"
    ${BINARY} add-genesis-account "${USER_ADDR}" "${COMMUNITY_TOKENS}" --home "${NODE_HOME}"
}

# Setup genesis transactions
setup_genesis_txs() {
    log "Setting up genesis transactions..."
    ${BINARY} gentx "${KEY_NAME}" "${STAKE_AMOUNT}" \
        --chain-id "${CHAIN_ID}" \
        --moniker "${MONIKER}" \
        --commission-rate "0.10" \
        --commission-max-rate "0.20" \
        --commission-max-change-rate "0.01" \
        --min-self-delegation "1000000" \
        --home "${NODE_HOME}" \
        ${KEYRING_FLAG}

    ${BINARY} collect-gentxs --home "${NODE_HOME}"
}

# Start the node
start_node() {
    log "Starting node..."
    ${BINARY} start \
        --home "${NODE_HOME}" \
        --pruning "nothing" \
        --minimum-gas-prices "${MIN_GAS_PRICES}" &

    # Give the node a moment to start
    sleep 2
    log "Waiting for node to start..."
    while ! curl -s http://localhost:26657/status > /dev/null; do
        log "Node not ready yet, waiting..."
        sleep 2
    done
    log "Node is up!"
    sleep 5
}

# Verify deployment
verify_deployment() {
    log "Verifying deployment..."
    
    # Check validator status
    VALIDATOR_ADDR=$(${BINARY} keys show "${KEY_NAME}" --bech val -a ${KEYRING_FLAG})
    VALIDATOR_STATUS=$(${BINARY} query staking validator "${VALIDATOR_ADDR}" \
        --chain-id "${CHAIN_ID}" \
        --node http://localhost:26657 \
        --output json)
    
    # Check user balance
    USER_BALANCE=$(${BINARY} query bank balances $(${BINARY} keys show "${KEY_NAME2}" -a ${KEYRING_FLAG}) \
        --chain-id "${CHAIN_ID}" \
        --node http://localhost:26657 \
        --output json | jq -r '.balances[0].amount')

    # Remove denomination from COMMUNITY_TOKENS for comparison
    EXPECTED_BALANCE=${COMMUNITY_TOKENS%$TOKEN_DENOM}
    if [ -n "${VALIDATOR_STATUS}" ] && [ "${USER_BALANCE}" = "${EXPECTED_BALANCE}" ]; then
        log "✅ Deployment verified successfully!"
        return 0
    else
        log "❌ Deployment verification failed!"
        return 1
    fi
}

# Main execution
main() {
    log "Starting deployment test for ${BINARY} blockchain..."
    cleanup
    create_keys
    init_chain
    setup_genesis_accounts
    setup_genesis_txs
    start_node
    verify_deployment
}

# Execute main function
main
