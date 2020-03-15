package generate

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	sdkAccounts "github.com/SebastianJ/harmony-sdk/accounts"
	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	sdkValidator "github.com/SebastianJ/harmony-sdk/staking/validator"
	sdkTxs "github.com/SebastianJ/harmony-sdk/transactions"
	"github.com/SebastianJ/harmony-stress/staking/delegation"
	"github.com/SebastianJ/harmony-stress/utils"
	"github.com/SebastianJ/harmony-tf/accounts"
	"github.com/SebastianJ/harmony-tf/balances"
	"github.com/SebastianJ/harmony-tf/funding"
	"github.com/SebastianJ/harmony-tf/logger"
	"github.com/SebastianJ/harmony-tf/testing"
	goSdkAccount "github.com/harmony-one/go-sdk/pkg/account"
	"github.com/harmony-one/harmony/numeric"
)

var (
	validators []string
	pools      = 1
)

// StressTestDelegations - stress tests delegations (generating accounts)
func StressTestDelegations() (err error) {
	fmt.Printf("Starting delegation stress tests - network: %s, mode: %s, node: %s\n", delegation.Configuration.Network.Name, delegation.Configuration.Network.Mode, delegation.Configuration.Network.Node)

	currentNonce := sdkNetwork.CurrentNonce(delegation.Configuration.Network.RPC, delegation.Configuration.Funding.Account.Address)
	gasPrice := delegation.Configuration.Network.Gas.Price

	validators, err = sdkValidator.All(delegation.Configuration.Network.RPC)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Found a total of %d total validators to send delegations to", len(validators)))

	if delegation.Configuration.Application.Count > delegation.Configuration.Application.PoolSize {
		pools = int(math.RoundToEven(float64(delegation.Configuration.Application.Count) / float64(delegation.Configuration.Application.PoolSize)))
	}

	for poolIndex := 0; poolIndex < pools; poolIndex++ {
		if poolIndex > 1 {
			currentNonce = sdkNetwork.CurrentNonce(delegation.Configuration.Network.RPC, delegation.Configuration.Funding.Account.Address)
			/*gasPrice = sdkTxs.BumpGasPrice(gasPrice)
			if gasPrice.GT(delegation.Configuration.Delegation.Gas.Price.Mul(numeric.NewDec(10))) {
				gasPrice = delegation.Configuration.Delegation.Gas.Price
			}*/

			fmt.Println(fmt.Sprintf("Nonce is now: %d, gas price is now: %f", currentNonce, gasPrice))
		}

		PerformDelegations(currentNonce, gasPrice)
	}

	return nil
}

// PerformDelegations - performs the actual account creation and delegation sending via goroutines
func PerformDelegations(nonce uint64, gasPrice numeric.Dec) {
	var waitGroup sync.WaitGroup

	for i := 0; i < delegation.Configuration.Application.PoolSize; i++ {
		waitGroup.Add(1)

		r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		validatorAddress := utils.RandomStringSliceItem(r, validators)

		go CreateAccountAndDelegate(validatorAddress, i, nonce, gasPrice, &waitGroup)
		nonce++
	}

	waitGroup.Wait()
}

// CreateAccountAndDelegate - creates a delegator account and delegates from it to the specified validator
func CreateAccountAndDelegate(validatorAddress string, index int, fundingNonce uint64, fundingGasPrice numeric.Dec, waitGroup *sync.WaitGroup) error {
	defer waitGroup.Done()
	accountName := fmt.Sprintf("CreateDelegation_%d_Account_%d", index, time.Now().UTC().UnixNano())
	logger.AccountLog(fmt.Sprintf("Generating a new account: %s", accountName), delegation.Configuration.Application.Verbose)

	account, err := accounts.GenerateAccount(accountName)
	if err != nil {
		logger.ErrorLog(err.Error(), delegation.Configuration.Application.Verbose)
		return err
	}

	logger.AccountLog(fmt.Sprintf("Generated account: %s, address: %s", account.Name, account.Address), delegation.Configuration.Application.Verbose)

	fundingAccountBalance, err := balances.GetShardBalance(delegation.Configuration.Funding.Account.Address, 0)
	if err != nil {
		logger.ErrorLog(fmt.Sprintf("Failed to retrieve shard balance - error: %s", err.Error()), delegation.Configuration.Application.Verbose)
		return err
	}

	logger.FundingLog(fmt.Sprintf("Available funding amount in the funding account %s, address: %s is %f", delegation.Configuration.Funding.Account.Name, delegation.Configuration.Funding.Account.Address, fundingAccountBalance), delegation.Configuration.Application.Verbose)
	fundingAmount := funding.CalculateFundingAmount(delegation.Configuration.Delegation.Amount, fundingAccountBalance, 1)
	logger.FundingLog(fmt.Sprintf("Funding account %s, address: %s with %f", account.Name, account.Address, fundingAmount), delegation.Configuration.Application.Verbose)
	funding.PerformFundingTransaction(&delegation.Configuration.Funding.Account, 0, account.Address, 0, fundingAmount, int(fundingNonce), delegation.Configuration.Funding.Gas.Limit, fundingGasPrice, delegation.Configuration.Funding.Timeout, delegation.Configuration.Funding.Attempts)

	accountStartingBalance, _ := balances.GetShardBalance(account.Address, 0)
	logger.AccountLog(fmt.Sprintf("Using account %s, address: %s to send a new delegation", account.Name, account.Address), delegation.Configuration.Application.Verbose)
	logger.BalanceLog(fmt.Sprintf("Account %s, address: %s has a starting balance of %f in shard %d before the test", account.Name, account.Address, accountStartingBalance, 0), delegation.Configuration.Application.Verbose)

	if !accountStartingBalance.IsZero() {
		delegation.Configuration.Account.Account = account

		logger.TransactionLog(fmt.Sprintf("Sending delegation transaction - will wait %d seconds for it to finalize", delegation.Configuration.Application.Timeout), delegation.Configuration.Application.Verbose)
		fmt.Println("")

		rawTx, err := executeDelegation(&account, validatorAddress, 0, delegation.Configuration.Delegation.Gas.Price)
		if err != nil {
			return nil
		}

		tx := sdkTxs.ToTransaction(account.Address, 0, account.Address, 0, rawTx, err)

		if delegation.Configuration.Application.Timeout > 0 {
			txResultColoring := logger.ResultColoring(tx.Success, true).Render(fmt.Sprintf("tx successful: %t", tx.Success))
			logger.TransactionLog(fmt.Sprintf("Performed account generation and delegation - delegator address: %s, validator address: %s, transaction hash: %s, %s", account.Address, validatorAddress, tx.TransactionHash, txResultColoring), delegation.Configuration.Application.Verbose)
		} else {
			logger.TransactionLog(fmt.Sprintf("Performed account generation and delegation - delegator address: %s, validator address: %s, transaction hash: %s", account.Address, validatorAddress, tx.TransactionHash), delegation.Configuration.Application.Verbose)
		}

		logger.TeardownLog("Performing test teardown (returning funds and removing accounts)\n", delegation.Configuration.Application.Verbose)

		if delegation.Configuration.Application.Timeout > 0 {
			testing.Teardown(&account, 0, delegation.Configuration.Funding.Account.Address, 0)
		} else {
			goSdkAccount.RemoveAccount(account.Name)
		}
	} else {
		logger.ErrorLog(fmt.Sprintf("Account %s, address: %s doesn't have sufficient balance to delegate to a validator! The account should have a balance of %f but the actual balance is %f", account.Name, account.Address, fundingAmount, accountStartingBalance), delegation.Configuration.Application.Verbose)
	}

	return nil
}

func executeDelegation(account *sdkAccounts.Account, address string, nonce uint64, gasPrice numeric.Dec) (map[string]interface{}, error) {
	txResult, err := delegation.Delegate(account, address, nonce, gasPrice)
	if err == nil {
		txHash := txResult["transactionHash"].(string)
		logger.TransactionLog(fmt.Sprintf("Sent delegation of %f from %s to %s, nonce: %d, gas price: %f, tx hash: %s", delegation.Configuration.Delegation.Amount, delegation.Configuration.Account.Account.Address, address, nonce, gasPrice, txHash), delegation.Configuration.Application.Verbose)
	} else {
		logger.ErrorLog(fmt.Sprintf("Error occurred while sending delegation of %f from %s to %s, nonce: %d - error: %s", delegation.Configuration.Delegation.Amount, delegation.Configuration.Account.Account.Address, address, nonce, err.Error()), delegation.Configuration.Application.Verbose)
	}

	return txResult, err
}
