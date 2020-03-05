package validators

import (
	"errors"
	"path/filepath"
	"strings"

	sdkAccounts "github.com/SebastianJ/harmony-sdk/accounts"
	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	cmd "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/SebastianJ/harmony-stress/utils"
	tfConfig "github.com/SebastianJ/harmony-tf/config"
	tfParams "github.com/SebastianJ/harmony-tf/testing/parameters"
	tfUtils "github.com/SebastianJ/harmony-tf/utils"
	"github.com/harmony-one/go-sdk/pkg/sharding"
)

// Configuration - the central configuration for the test suite tool
var Configuration Config

// Staking - staking and validator settings used for creating validators
var Staking tfParams.StakingParameters

// Configure - configures the test suite tool using a combination of the YAML config file as well as command arguments
func Configure(basePath string, flags cmd.PersistentFlags) (err error) {
	configPath := filepath.Join(basePath, "config.yml")
	if err = loadYamlConfig(configPath); err != nil {
		return err
	}

	stakingPath := filepath.Join(basePath, "staking.yml")
	if err = loadStakingConfig(stakingPath); err != nil {
		return err
	}

	if Configuration.BasePath == "" {
		Configuration.BasePath = basePath
	}

	Configuration.Application.Verbose = true // this configures the output using Harmony TF's logger
	sdkNetwork.Verbose = flags.Verbose       // this configures the raw tx dump logs from Harmony-SDK

	if err = configureNetworkConfig(flags); err != nil {
		return err
	}

	if err = configureBaseConfig(flags); err != nil {
		return err
	}

	tfConfig.ConfigureStylingConfig()

	return nil
}

func configureNetworkConfig(flags cmd.PersistentFlags) (err error) {
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

	if flags.Node != "" && flags.Node != Configuration.Network.Node {
		Configuration.Network.Node = flags.Node
	}

	Configuration.Network.API = sdkNetwork.Network{
		Name: Configuration.Network.Name,
		Mode: Configuration.Network.Mode,
		Node: Configuration.Network.Node,
	}

	Configuration.Network.Gas.Initialize()
	Configuration.Network.API.Initialize()

	Configuration.Network.Node = Configuration.Network.API.NodeAddress(0)

	shardingStructure, err := sharding.Structure(Configuration.Network.Node)
	if err != nil {
		return err
	}

	Configuration.Network.Shards = len(shardingStructure)

	Configuration.Network.RPC, err = Configuration.Network.API.RPCClient(0)
	if err != nil {
		return err
	}

	tfConfig.Configuration.Network = tfConfig.Network{
		Name:   Configuration.Network.Name,
		Mode:   Configuration.Network.Mode,
		Node:   Configuration.Network.Node,
		Shards: Configuration.Network.Shards,
		API:    Configuration.Network.API,
		Gas:    Configuration.Network.Gas,
	}

	return nil
}

func configureBaseConfig(flags cmd.PersistentFlags) error {
	fromAddress := flags.From
	if fromAddress == "" {
		return errors.New("you need to specify the sender address")
	}

	Configuration.Funding.Account.Address = fromAddress

	if Configuration.Funding.Account.Name == "" {
		Configuration.Funding.Account.Name = sdkAccounts.FindAccountNameByAddress(Configuration.Funding.Account.Address)
	}

	if flags.Passphrase != "" && flags.Passphrase != Configuration.Application.Passphrase {
		Configuration.Application.Passphrase = flags.Passphrase
	}

	Configuration.Funding.Account.Passphrase = Configuration.Application.Passphrase
	tfConfig.Configuration.Account.Passphrase = Configuration.Application.Passphrase

	Configuration.Funding.Gas.Initialize()

	tfConfig.Configuration.Funding = tfConfig.Funding{
		Account:  Configuration.Funding.Account,
		Timeout:  Configuration.Funding.Timeout,
		Attempts: Configuration.Funding.Attempts,
		Gas:      Configuration.Funding.Gas,
	}

	tfConfig.Configuration.Funding.Timeout = tfUtils.NetworkTimeoutAdjustment(tfConfig.Configuration.Network.Name, tfConfig.Configuration.Funding.Timeout)

	Configuration.Application.Infinite = flags.Infinite

	if flags.Count > 0 && flags.Count != Configuration.Application.Count {
		Configuration.Application.Count = flags.Count
	}

	if flags.PoolSize > 0 && flags.PoolSize != Configuration.Application.PoolSize {
		Configuration.Application.PoolSize = flags.PoolSize
	}

	// Initialize the staking params - this converts values to numeric.Dec etc.
	Staking.Initialize()
	Staking.Timeout = tfUtils.NetworkTimeoutAdjustment(tfConfig.Configuration.Network.Name, Staking.Timeout)

	return nil
}

func loadYamlConfig(path string) error {
	Configuration = Config{}
	return utils.ParseYaml(path, &Configuration)
}

func loadStakingConfig(path string) error {
	Staking = tfParams.StakingParameters{}
	return utils.ParseYaml(path, &Staking)
}
