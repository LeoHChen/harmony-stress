package transactions

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	cmd "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/SebastianJ/harmony-stress/utils"
	goSdkCommon "github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/go-sdk/pkg/store"
	"gopkg.in/yaml.v2"
)

// Configuration - the central configuration for the test suite tool
var Configuration Config

// Configure - configures the test suite tool using a combination of the YAML config file as well as command arguments
func Configure(basePath string, flags cmd.PersistentFlags, txFlags cmd.TxFlags) (err error) {
	configPath := filepath.Join(basePath, "config.yml")
	if err = loadYamlConfig(configPath); err != nil {
		return err
	}

	if Configuration.BasePath == "" {
		Configuration.BasePath = basePath
	}

	Configuration.Verbose = flags.Verbose
	// Set the verbosity level of harmony-sdk
	sdkNetwork.Verbose = Configuration.Verbose

	goSdkCommon.DebugRPC = flags.VerboseGoSdk

	// It's very important that configureTransactionsConfig gets executed first since it sets config fields that are later used by other configuration steps
	if err = configureTransactionsConfig(flags, txFlags); err != nil {
		return err
	}

	if err = configureNetworkConfig(flags, txFlags); err != nil {
		return err
	}

	if err = configureAccountConfig(flags, txFlags); err != nil {
		return err
	}

	return nil
}

func configureTransactionsConfig(flags cmd.PersistentFlags, txFlags cmd.TxFlags) error {
	if txFlags.Mode != "" && txFlags.Mode != Configuration.Transactions.Mode {
		Configuration.Transactions.Mode = txFlags.Mode
	}

	if txFlags.FromShardID >= 0 && uint32(txFlags.FromShardID) != Configuration.Transactions.FromShard {
		Configuration.Transactions.FromShard = uint32(txFlags.FromShardID)
	}

	if txFlags.ToShardID >= 0 && uint32(txFlags.ToShardID) != Configuration.Transactions.ToShard {
		Configuration.Transactions.ToShard = uint32(txFlags.ToShardID)
	}

	fromAddress := flags.From
	if fromAddress == "" {
		return errors.New("you need to specify the sender address")
	}
	Configuration.Transactions.From = fromAddress

	txDataFilePath := filepath.Join(Configuration.BasePath, "data/data.txt")
	txData, err := utils.ReadFileToString(txDataFilePath)
	if err == nil && txData != "" && txData != Configuration.Transactions.Data {
		Configuration.Transactions.Data = txData
	}

	if txFlags.Amount != "" && txFlags.Amount != Configuration.Transactions.RawAmount {
		Configuration.Transactions.RawAmount = txFlags.Amount
	}

	if flags.GasPrice != "" && flags.GasPrice != Configuration.Transactions.Gas.RawPrice {
		Configuration.Transactions.Gas.RawPrice = flags.GasPrice
	}

	if flags.Count >= 0 && flags.Count != Configuration.Transactions.Count {
		Configuration.Transactions.Count = flags.Count
	}

	if flags.PoolSize >= 0 && flags.PoolSize != Configuration.Transactions.PoolSize {
		Configuration.Transactions.PoolSize = flags.PoolSize
	}

	if flags.Timeout >= 0 && flags.Timeout != Configuration.Transactions.Timeout {
		Configuration.Transactions.Timeout = flags.Timeout
	}

	Configuration.Transactions.Initialize()

	receiversPath := filepath.Join(Configuration.BasePath, "data/receivers.txt")
	receivers, _ := FetchReceivers(receiversPath)

	if len(receivers) == 0 {
		return fmt.Errorf("you need to create the file %s and add at least one receiver address to it", receiversPath)
	}

	Configuration.Transactions.Receivers = receivers

	return nil
}

func configureNetworkConfig(flags cmd.PersistentFlags, txFlags cmd.TxFlags) (err error) {
	if flags.Network != "" && flags.Network != Configuration.Network.Name {
		Configuration.Network.Name = flags.Network
	}

	Configuration.Network.Name = sdkNetwork.NormalizedNetworkName(Configuration.Network.Name)
	if Configuration.Network.Name == "" {
		return errors.New("you need to specify a valid network name to use! Valid options: localnet, devnet, testnet, staking or mainnet")
	}

	Configuration.Network.Mode = strings.ToLower(Configuration.Network.Mode)
	mode := strings.ToLower(flags.NetworkMode)
	if mode != "" && mode != Configuration.Network.Mode {
		Configuration.Network.Mode = mode
	}

	Configuration.Network.API = sdkNetwork.Network{
		Name: Configuration.Network.Name,
		Mode: Configuration.Network.Mode,
	}

	Configuration.Network.API.Initialize()
	// Temporarily hard code nodes to work around RPC limits

	if flags.Node != "" && flags.Node != Configuration.Network.Node {
		Configuration.Network.Node = flags.Node
	}

	if Configuration.Network.Node == "" {
		Configuration.Network.Node = Configuration.Network.API.NodeAddress(Configuration.Transactions.FromShard)
	}

	if Configuration.Verbose {
		fmt.Printf("Using network: %s, mode: %s, node: %s\n", Configuration.Network.Name, Configuration.Network.Mode, Configuration.Network.Node)
	}

	shardingStructure, err := sharding.Structure(Configuration.Network.Node)
	if err != nil {
		return err
	}

	Configuration.Network.Shards = len(shardingStructure)

	Configuration.Network.RPC, err = sdkNetwork.NewRPCClient(Configuration.Network.Node, Configuration.Transactions.FromShard)
	if err != nil {
		return err
	}

	return nil
}

func configureAccountConfig(flags cmd.PersistentFlags, txFlags cmd.TxFlags) (err error) {
	if flags.Passphrase != Configuration.Account.Passphrase {
		Configuration.Account.Passphrase = flags.Passphrase
	}

	Configuration.Account.Keystore, Configuration.Account.Account, err = store.UnlockedKeystore(Configuration.Transactions.From, Configuration.Account.Passphrase)
	if err != nil {
		return err
	}

	Configuration.Account.Nonce = sdkNetwork.CurrentNonce(Configuration.Network.RPC, Configuration.Transactions.From)

	return nil
}

func loadYamlConfig(path string) error {
	Configuration = Config{}

	yamlData, err := utils.ReadFileToString(path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(yamlData), &Configuration)

	if err != nil {
		return err
	}

	return nil
}
