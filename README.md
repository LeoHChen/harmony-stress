# Harmony Stress
Stress testing tools for Harmony.

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

### Regular transaction stress testing

```
./stress txs --from YOUR_ADDRESS --network NETWORK --count COUNT --pool-size POOL_SIZE
```

### Validator creation stress testing

```
./stress validators --from YOUR_ADDRESS --network NETWORK --count COUNT --pool-size POOL_SIZE
```

### Delegation stress testing

```
dist/stress delegations --from YOUR_ADDRESS --network NETWORK --count COUNT --pool-size POOL_SIZE
```

### All options:

```
$ ./stress --help
Harmony stress test tool

Usage:
  stress [flags]
  stress [command]

Available Commands:
  delegations Stress test delegation transactions
  help        Help about any command
  txs         Stress test normal transactions
  validators  Stress test validator creation
  version     Show version

Flags:
      --app-mode string       <app-mode> (default "async")
      --count int             <count> (default 1000)
      --from string           <from>
  -h, --help                  help for stress
      --infinite              <infinite>
      --network string        <network> (default "localnet")
      --network-mode string   <mode> (default "api")
      --node string           <node>
      --passphrase string     <passphrase>
      --path string           <path> (default ".")
      --pool-size int         <pool-size> (default 100)
      --pprof-port int        <pprof-port> (default -1)
      --timeout int           <pool-size>
      --verbose               <verbose>
      --verbose-go-sdk        <verbose-go-sdk>

Use "stress [command] --help" for more information about a command.
```
