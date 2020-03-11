package delegation

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	sdkAccounts "github.com/SebastianJ/harmony-sdk/accounts"
	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	cmd "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/SebastianJ/harmony-stress/utils"
	goSdkCommon "github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"gopkg.in/yaml.v2"
)

// Configuration - the central configuration for the test suite tool
var Configuration Config

// Configure - configures the test suite tool using a combination of the YAML config file as well as command arguments
func Configure(basePath string, flags cmd.PersistentFlags, delegationFlags cmd.DelegationFlags) (err error) {
	configPath := filepath.Join(basePath, "config.yml")
	if err = loadYamlConfig(configPath); err != nil {
		return err
	}

	if Configuration.BasePath == "" {
		Configuration.BasePath = basePath
	}

	Configuration.Application.Verbose = flags.Verbose
	sdkNetwork.Verbose = Configuration.Application.Verbose // Set the verbosity level of harmony-sdk
	goSdkCommon.DebugRPC = flags.VerboseGoSdk

	// It's very important that configureTransactionsConfig gets executed first since it sets config fields that are later used by other configuration steps
	if err = configureDelegationConfig(flags, delegationFlags); err != nil {
		return err
	}

	if err = configureNetworkConfig(flags, delegationFlags); err != nil {
		return err
	}

	if err = configureAccountConfig(flags, delegationFlags); err != nil {
		return err
	}

	return nil
}

func configureDelegationConfig(flags cmd.PersistentFlags, delegationFlags cmd.DelegationFlags) error {
	if flags.ApplicationMode != "" && flags.ApplicationMode != Configuration.Application.Mode {
		Configuration.Application.Mode = flags.ApplicationMode
	}

	fromAddress := flags.From
	if fromAddress == "" {
		return errors.New("you need to specify the sender address")
	}
	Configuration.Application.From = fromAddress

	if flags.Count >= 0 && flags.Count != Configuration.Application.Count {
		Configuration.Application.Count = flags.Count
	}

	if flags.PoolSize >= 0 && flags.PoolSize != Configuration.Application.PoolSize {
		Configuration.Application.PoolSize = flags.PoolSize
	}

	if flags.Timeout >= 0 && flags.Timeout != Configuration.Application.Timeout {
		Configuration.Application.Timeout = flags.Timeout
	}

	if delegationFlags.FromShardID >= 0 && uint32(delegationFlags.FromShardID) != Configuration.Delegation.FromShard {
		Configuration.Delegation.FromShard = uint32(delegationFlags.FromShardID)
	}

	if delegationFlags.ToShardID >= 0 && uint32(delegationFlags.ToShardID) != Configuration.Delegation.ToShard {
		Configuration.Delegation.ToShard = uint32(delegationFlags.ToShardID)
	}

	if delegationFlags.Amount != "" && delegationFlags.Amount != Configuration.Delegation.RawAmount {
		Configuration.Delegation.RawAmount = delegationFlags.Amount
	}

	if flags.GasPrice != "" && flags.GasPrice != Configuration.Delegation.Gas.RawPrice {
		Configuration.Delegation.Gas.RawPrice = flags.GasPrice
	}

	Configuration.Delegation.Initialize()

	return nil
}

func configureNetworkConfig(flags cmd.PersistentFlags, delegationFlags cmd.DelegationFlags) (err error) {
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
		Configuration.Network.Node = Configuration.Network.API.NodeAddress(Configuration.Delegation.FromShard)
	}

	if Configuration.Application.Verbose {
		fmt.Printf("Using network: %s, mode: %s, node: %s\n", Configuration.Network.Name, Configuration.Network.Mode, Configuration.Network.Node)
	}

	shardingStructure, err := sharding.Structure(Configuration.Network.Node)
	if err != nil {
		return err
	}

	Configuration.Network.Shards = len(shardingStructure)

	Configuration.Network.RPC, err = sdkNetwork.NewRPCClient(Configuration.Network.Node, Configuration.Delegation.FromShard)
	if err != nil {
		return err
	}

	return nil
}

func configureAccountConfig(flags cmd.PersistentFlags, delegationFlags cmd.DelegationFlags) (err error) {
	Configuration.Account.Account = sdkAccounts.Account{
		Address: Configuration.Application.From,
	}

	if flags.Passphrase != Configuration.Account.Passphrase {
		Configuration.Account.Account.Passphrase = flags.Passphrase
	}

	Configuration.Account.Account.Unlock()

	Configuration.Account.Nonce = sdkNetwork.CurrentNonce(Configuration.Network.RPC, Configuration.Application.From)

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
