package delegation

import (
	sdkAccounts "github.com/SebastianJ/harmony-sdk/accounts"
	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	"github.com/harmony-one/go-sdk/pkg/common"
	goSdkRpc "github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/numeric"
)

// Config - represents the config
type Config struct {
	BasePath    string      `yaml:"-"`
	Network     Network     `yaml:"network"`
	Application Application `yaml:"application"`
	Delegation  Delegation  `yaml:"delegation"`
	Account     Account     `yaml:"account"`
}

// Application - represents the transactions settings group
type Application struct {
	From       string
	Mode       string `yaml:"mode"`
	Infinite   bool   `yaml:"infinite"`
	Count      int    `yaml:"count"`
	PoolSize   int    `yaml:"pool_size"`
	Verbose    bool   `yaml:"verbose"`
	Passphrase string `yaml:"passphrase"`
	Timeout    int    `yaml:"timeout"`
}

// Delegation - represents the transactions settings group
type Delegation struct {
	RawAmount string         `yaml:"amount"`
	Amount    numeric.Dec    `yaml:"-"`
	FromShard uint32         `yaml:"from_shard"`
	ToShard   uint32         `yaml:"to_shard"`
	Gas       sdkNetwork.Gas `yaml:"gas"`
}

// Account - represents the account settings group
type Account struct {
	Passphrase string `yaml:"passphrase"`
	Nonce      uint64
	Account    sdkAccounts.Account
}

// Network - represents the network settings group
type Network struct {
	Name   string                  `yaml:"name"`
	Mode   string                  `yaml:"mode"`
	Node   string                  `yaml:"-"`
	Shards int                     `yaml:"-"`
	Gas    sdkNetwork.Gas          `yaml:"gas"`
	RPC    *goSdkRpc.HTTPMessenger `yaml:"-"`
	API    sdkNetwork.Network      `yaml:"-"`
}

// Initialize - convert values to the appropriate data types + ensure correct values
func (del *Delegation) Initialize() {
	if decAmount, err := common.NewDecFromString(del.RawAmount); err == nil {
		del.Amount = decAmount
	}

	del.Gas.Initialize()
}
