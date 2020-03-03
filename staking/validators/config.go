package validators

import (
	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	tfConfig "github.com/SebastianJ/harmony-tf/config"
	"github.com/gookit/color"
	goSdkRpc "github.com/harmony-one/go-sdk/pkg/rpc"
)

// Config - represents the config
type Config struct {
	BasePath    string           `yaml:"-"`
	Network     Network          `yaml:"network"`
	Application Application      `yaml:"application"`
	Funding     tfConfig.Funding `yaml:"funding"`
	Framework   Framework        `yaml:"framework"`
}

// Application - represents the transactions settings group
type Application struct {
	Infinite   bool   `yaml:"infinite"`
	Count      int    `yaml:"count"`
	PoolSize   int    `yaml:"pool_size"`
	Verbose    bool   `yaml:"verbose"`
	Passphrase string `yaml:"passphrase"`
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

// Framework - represents common framework settings for Harmony TF
type Framework struct {
	Styling Styling `yaml:"-"`
}

// Styling - represents settings for styling the log output
type Styling struct {
	Header         *color.Style `yaml:"-"`
	TestCaseHeader *color.Style `yaml:"-"`
	Default        *color.Style `yaml:"-"`
	Account        *color.Style `yaml:"-"`
	Funding        *color.Style `yaml:"-"`
	Balance        *color.Style `yaml:"-"`
	Transaction    *color.Style `yaml:"-"`
	Staking        *color.Style `yaml:"-"`
	Teardown       *color.Style `yaml:"-"`
	Success        *color.Style `yaml:"-"`
	Error          *color.Style `yaml:"-"`
	Padding        string       `yaml:"-"`
}
