package transactions

import (
	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/numeric"
)

// Config - represents the config
type Config struct {
	BasePath     string       `yaml:"-"`
	Verbose      bool         `yaml:"-"`
	Transactions Transactions `yaml:"transactions"`
	Network      Network      `yaml:"network"`
	Account      Account      `yaml:"account"`
}

// Transactions - represents the transactions settings group
type Transactions struct {
	Mode      string `yaml:"mode"`
	From      string
	RawAmount string         `yaml:"amount"`
	Amount    numeric.Dec    `yaml:"-"`
	FromShard uint32         `yaml:"from_shard"`
	ToShard   uint32         `yaml:"to_shard"`
	Data      string         `yaml:"data"`
	Count     int            `yaml:"count"`
	PoolSize  int            `yaml:"pool_size"`
	Timeout   int            `yaml:"timeout"`
	Gas       sdkNetwork.Gas `yaml:"gas"`
	Receivers []string
}

// Network - represents the network settings group
type Network struct {
	Name   string             `yaml:"name"`
	Mode   string             `yaml:"mode"`
	Node   string             `yaml:"-"`
	Shards int                `yaml:"-"`
	RPC    *rpc.HTTPMessenger `yaml:"-"`
	API    sdkNetwork.Network `yaml:"-"`
}

// Account - represents the account settings group
type Account struct {
	Passphrase string `yaml:"passphrase"`
	Nonce      uint64
	Keystore   *keystore.KeyStore
	Account    *accounts.Account
}

// Initialize - convert values to the appropriate data types + ensure correct values
func (txs *Transactions) Initialize() {
	if decAmount, err := common.NewDecFromString(txs.RawAmount); err == nil {
		txs.Amount = decAmount
	}

	txs.Gas.Initialize()
}
