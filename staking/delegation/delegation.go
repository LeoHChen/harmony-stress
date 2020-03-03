package delegation

import (
	"errors"

	sdkDelegation "github.com/SebastianJ/harmony-sdk/staking/delegation"
)

var (
	errNilDelegate   = errors.New("Delegate can not be nil or a negative value")
	errNilUndelegate = errors.New("Undelegate can not be nil or a negative value")
)

// Delegate - performs delegation
func Delegate(validatorAddress string, nonce uint64) (map[string]interface{}, error) {
	return executeDelegationMethod("delegate", validatorAddress, nonce)
}

// Undelegate - performs undelegation
func Undelegate(validatorAddress string, nonce uint64) (map[string]interface{}, error) {
	return executeDelegationMethod("undelegate", validatorAddress, nonce)
}

func executeDelegationMethod(method string, validatorAddress string, nonce uint64) (txResult map[string]interface{}, err error) {
	Configuration.Account.Account.Unlock()

	if method == "delegate" {
		txResult, err = sdkDelegation.Delegate(
			Configuration.Account.Account.Keystore,
			Configuration.Account.Account.Account,
			Configuration.Network.RPC,
			Configuration.Network.API.ChainID,
			Configuration.Account.Account.Address,
			validatorAddress,
			Configuration.Delegation.Amount,
			Configuration.Delegation.Gas.Limit,
			Configuration.Delegation.Gas.Price,
			nonce,
			Configuration.Account.Account.Passphrase,
			Configuration.Network.API.NodeAddress(0),
			Configuration.Application.Timeout,
		)
	} else if method == "undelegate" {
		txResult, err = sdkDelegation.Undelegate(
			Configuration.Account.Account.Keystore,
			Configuration.Account.Account.Account,
			Configuration.Network.RPC,
			Configuration.Network.API.ChainID,
			Configuration.Account.Account.Address,
			validatorAddress,
			Configuration.Delegation.Amount,
			Configuration.Delegation.Gas.Limit,
			Configuration.Delegation.Gas.Price,
			nonce,
			Configuration.Account.Account.Passphrase,
			Configuration.Network.API.NodeAddress(0),
			Configuration.Application.Timeout,
		)
	}

	if err != nil {
		return nil, err
	}

	return txResult, nil
}
