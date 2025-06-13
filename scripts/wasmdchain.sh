#!/usr/bin/env bash

export GENESIS_ACCOUNT="wasmaccount"
export FAUCET_ACCOUNT="faucetaccount"

export CHAIN_ID=${CHAIN_ID:-"wasm-testnet"}
export MONIKER="wasmchain"
export BINARY=${BINARY:-wasmd}
export KEYALGO="secp256k1"
export KEYRING=${KEYRING:-"test"}
export DENOM=${DENOM:-uatom}
export HOME_DIR=$(eval echo "${HOME_DIR:-"~/.wasmd"}")

export CLEAN=${CLEAN:-"true"}
export RPC=${RPC:-"26657"}
export REST=${REST:-"1317"}
export GRPC=${GRPC:-"9090"}
export PROFF=${PROFF:-"6060"}
export P2P=${P2P:-"26656"}
export GRPC_WEB=${GRPC_WEB:-"9091"}
export ROSETTA=${ROSETTA:-"8080"}
export BLOCK_TIME=${BLOCK_TIME:-"5s"}

# wasm1pdgy8myz7pwvqwwvnjz5nws9vfu5uchzw0543p
export GENESIS_MNEMONIC=${GENESIS_MNEMONIC:-"dress volcano dwarf school nature sing security domain lady video attitude jaguar miss unit detect title spot embark apple style cotton trouble pudding west"};
# wasm1g6ckzglt6uxc9nf7m2fpveym9vp3x8mxvvdhrh
export FAUCET_MNEMONIC=${FAUCET_MNEMONIC:-"detect outdoor science mind peace marriage decrease stage head kangaroo tuna enforce chat video goose minute warfare creek advance zebra define mesh final cheap"};

# If which binary does not exist, exit
if [ -z `which $BINARY` ]; then
  echo "Ensure $BINARY is installed and in your PATH"
  exit
fi

alias BINARY="$BINARY --home=$HOME_DIR"

command -v $BINARY > /dev/null 2>&1 || { echo >&2 "$BINARY command not found. Ensure this is setup / properly installed in your GOPATH (make install)."; exit 1; }
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

set_config() {
  $BINARY config set client chain-id $CHAIN_ID
  $BINARY config set client keyring-backend $KEYRING
  $BINARY config set client node "tcp://localhost:$RPC"
}

from_chain_setup () {

  # remove existing daemon files.
  if [ ${#HOME_DIR} -le 2 ]; then
      echo "HOME_DIR must be more than 2 characters long"
      return
  fi
  rm -rf $HOME_DIR && echo "Removed $HOME_DIR"

  # reset values if not set already after whipe
  set_config

  # add accounts
  add_key() {
    key=$1
    mnemonic=$2
    if [ "$KEYRING" = "test" ]; then
        echo $mnemonic | $BINARY keys add $key --keyring-backend $KEYRING --recover --home=$HOME_DIR
    elif [ "$KEYRING" = "os" ] || [ "$KEYRING" = "file" ]; then
        $BINARY keys add $key --keyring-backend $KEYRING --recover --home=$HOME_DIR
    else
      echo "Sorry, This script does not support keying's backend being of type 'kwallet | pass | memory'"
    fi
  }

  # add genesis and faucet account.
  add_key "$GENESIS_ACCOUNT" "$GENESIS_MNEMONIC"
  add_key "$FAUCET_ACCOUNT" "$FAUCET_MNEMONIC"

  # chain initial setup
  $BINARY init $MONIKER --chain-id $CHAIN_ID

  update_genesis () {
    cat $HOME_DIR/config/genesis.json | jq "$1" > $HOME_DIR/config/tmp_genesis.json && mv $HOME_DIR/config/tmp_genesis.json $HOME_DIR/config/genesis.json
  }

  # === CORE MODULES ===

  # Block
  update_genesis '.consensus_params["block"]["max_gas"]="100000000"'

  # staking
  update_genesis `printf '.app_state["staking"]["params"]["bond_denom"]="%s"' $DENOM`

  # mint
  update_genesis `printf '.app_state["mint"]["params"]["mint_denom"]="%s"' $DENOM`

  # crisis
  update_genesis `printf '.app_state["crisis"]["constant_fee"]={"denom":"%s","amount":"1000"}' $DENOM`

  # Allocate genesis accounts
  $BINARY genesis add-genesis-account $GENESIS_ACCOUNT 99900000000000$DENOM --keyring-backend $KEYRING
  $BINARY genesis add-genesis-account $FAUCET_ACCOUNT 100000000000$DENOM --keyring-backend $KEYRING

  # Sign genesis transaction
  $BINARY genesis gentx $GENESIS_ACCOUNT 10000000000000$DENOM --keyring-backend $KEYRING --chain-id $CHAIN_ID

  $BINARY genesis collect-gentxs

  $BINARY genesis validate-genesis
  err=$?
  if [ $err -ne 0 ]; then
    echo "Failed to validate genesis"
    return
  fi
}

# check if CLEAN is not set to false
if [ "$CLEAN" != "false" ]; then
  echo "Starting from a clean state"
  from_chain_setup
fi

echo "Starting wasmd node..."

os_type=$(uname -s)
# Linux
if [[ "$os_type" == "Linux" ]]; then
  # Opens the RPC endpoint to outside connections
  sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:'$RPC'"/g' $HOME_DIR/config/config.toml
  sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' $HOME_DIR/config/config.toml

  # Other config
  sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001uatom"/g' $HOME_DIR/config/app.toml
  
  # REST endpoint
  sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:'$REST'"/g' $HOME_DIR/config/app.toml
  sed -i 's/enable = false/enable = true/g' $HOME_DIR/config/app.toml
  sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME_DIR/config/app.toml

  # peer exchange
  sed -i 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:'$PROFF_LADDER'"/g' $HOME_DIR/config/config.toml
  sed -i 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:'$P2P'"/g' $HOME_DIR/config/config.toml

  # GRPC
  sed -i 's/address = "localhost:9090"/address = "0.0.0.0:'$GRPC'"/g' $HOME_DIR/config/app.toml
  sed -i 's/address = "localhost:9091"/address = "0.0.0.0:'$GRPC_WEB'"/g' $HOME_DIR/config/app.toml

  # Rosetta Api
  sed -i 's/address = ":8080"/address = "0.0.0.0:'$ROSETTA'"/g' $HOME_DIR/config/app.toml

  # Faster blocks
  sed -i 's/timeout_commit = "5s"/timeout_commit = "'$BLOCK_TIME'"/g' $HOME_DIR/config/config.toml

# Darwin
elif [[ "$os_type" == "Darwin" ]]; then
  sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:'$RPC'"/g' $HOME_DIR/config/config.toml
  sed -i '' 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' $HOME_DIR/config/config.toml

  sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001uatom"/g' $HOME_DIR/config/app.toml

  sed -i '' 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:'$REST'"/g' $HOME_DIR/config/app.toml
  sed -i '' 's/enable = false/enable = true/g' $HOME_DIR/config/app.toml
  sed -i '' 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME_DIR/config/app.toml

  sed -i '' 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:'$PROFF_LADDER'"/g' $HOME_DIR/config/config.toml
  sed -i '' 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:'$P2P'"/g' $HOME_DIR/config/config.toml

  sed -i '' 's/address = "localhost:9090"/address = "0.0.0.0:'$GRPC'"/g' $HOME_DIR/config/app.toml
  sed -i '' 's/address = "localhost:9091"/address = "0.0.0.0:'$GRPC_WEB'"/g' $HOME_DIR/config/app.toml

  sed -i '' 's/address = ":8080"/address = "0.0.0.0:'$ROSETTA'"/g' $HOME_DIR/config/app.toml

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "'$BLOCK_TIME'"/g' $HOME_DIR/config/config.toml

else
    echo "Unknown system type: $os_type"
    exit 1
fi

# Start the node
$BINARY start --pruning=nothing --rpc.laddr="tcp://0.0.0.0:$RPC"