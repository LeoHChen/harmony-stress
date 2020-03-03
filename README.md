# Harmony Validator Spammer
Harmony tx sender is a tool to bulk send transactions on Harmony's blockchain.

## Prerequisites
You need to import an existing key with funds to the keystore.

hmy is automatically downloaded as a part of the installation script.

Import a key using the following command:
```
./hmy keys import-ks ABSOLUTE_PATH_TO_YOUR_KEY NAME_OF_YOUR_KEY --passphrase ""
```

Find the address of your newly imported key:
```
./hmy keys list
```

## Installation

```
bash <(curl -s -S -L https://raw.githubusercontent.com/SebastianJ/harmony-stress/master/scripts/install.sh)
```

The installer script will also download the `config.yml` (contains general settings) and `staking.yml` (contains the create validator settings).


## Usage
```
./harmony-stress --network staking --from YOUR_SENDER_ACCOUNT_ADDRESS --infinite
```

### All options:

```
NAME:
   Harmony Validator Spammer - stress tests a staking enabled Harmony network/blockchain - Use --help to see all available arguments

USAGE:
   harmony-stress [global options] command [command options] [arguments...]

VERSION:
   go1.13.7/darwin-amd64

AUTHOR:
   Sebastian Johnsson

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --network value                 Which network to use (valid options: localnet, devnet, testnet, mainnet)
   --from value                    Which address to send tokens from (must exist in the keystore)
   --passphrase value              Passphrase to use for unlocking the keystore
   --infinite                      If the program should run in an infinite loop
   --count value                   How many transactions to send in total (default: 0)
   --pool-size value               How many validators to create simultaneously (default: 0)
   --timeout value  How long to wait for transactions to get confirmed (default: 0)
   --verbose                       Enable more verbose output
   --help, -h                      show help
   --version, -v                   print the version
```
