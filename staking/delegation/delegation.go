package delegation

import (
	"errors"

	sdkAccounts "github.com/SebastianJ/harmony-sdk/accounts"
	sdkDelegation "github.com/SebastianJ/harmony-sdk/staking/delegation"
	"github.com/harmony-one/harmony/numeric"
)

var (
	errNilDelegate   = errors.New("Delegate can not be nil or a negative value")
	errNilUndelegate = errors.New("Undelegate can not be nil or a negative value")
)

// Delegate - performs delegation
func Delegate(account *sdkAccounts.Account, validatorAddress string, nonce uint64, gasPrice numeric.Dec) (map[string]interface{}, error) {
	return executeDelegationMethod(account, "delegate", validatorAddress, nonce, gasPrice)
}

// Undelegate - performs undelegation
func Undelegate(account *sdkAccounts.Account, validatorAddress string, nonce uint64, gasPrice numeric.Dec) (map[string]interface{}, error) {
	return executeDelegationMethod(account, "undelegate", validatorAddress, nonce, gasPrice)
}

func executeDelegationMethod(account *sdkAccounts.Account, method string, validatorAddress string, nonce uint64, gasPrice numeric.Dec) (txResult map[string]interface{}, err error) {
	account.Unlock()

	if method == "delegate" {
		txResult, err = sdkDelegation.Delegate(
			account.Keystore,
			account.Account,
			Configuration.Network.RPC,
			Configuration.Network.API.ChainID,
			account.Address,
			validatorAddress,
			Configuration.Delegation.Amount,
			Configuration.Delegation.Gas.Limit,
			gasPrice,
			nonce,
			account.Passphrase,
			Configuration.Network.API.NodeAddress(0),
			Configuration.Application.Timeout,
		)
	} else if method == "undelegate" {
		txResult, err = sdkDelegation.Undelegate(
			account.Keystore,
			account.Account,
			Configuration.Network.RPC,
			Configuration.Network.API.ChainID,
			account.Address,
			validatorAddress,
			Configuration.Delegation.Amount,
			Configuration.Delegation.Gas.Limit,
			gasPrice,
			nonce,
			account.Passphrase,
			Configuration.Network.API.NodeAddress(0),
			Configuration.Application.Timeout,
		)
	}

	if err != nil {
		return nil, err
	}

	return txResult, nil
}
