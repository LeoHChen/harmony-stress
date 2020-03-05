package create

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sync"

	sdkAccounts "github.com/SebastianJ/harmony-sdk/accounts"
	sdkCrypto "github.com/SebastianJ/harmony-sdk/crypto"
	sdkNetwork "github.com/SebastianJ/harmony-sdk/network"
	sdkTxs "github.com/SebastianJ/harmony-sdk/transactions"
	"github.com/SebastianJ/harmony-stress/staking/validators"
	"github.com/SebastianJ/harmony-stress/utils"
	"github.com/SebastianJ/harmony-tf/accounts"
	"github.com/SebastianJ/harmony-tf/balances"
	"github.com/SebastianJ/harmony-tf/crypto"
	"github.com/SebastianJ/harmony-tf/funding"
	"github.com/SebastianJ/harmony-tf/logger"
	tfStaking "github.com/SebastianJ/harmony-tf/staking"
	"github.com/SebastianJ/harmony-tf/testing"
	goSdkAccount "github.com/harmony-one/go-sdk/pkg/account"
	"github.com/harmony-one/harmony/numeric"
)

// StressTestValidatorCreation - stress tests validator creation
func StressTestValidatorCreation() {
	fmt.Printf("Starting validator spammer - network: %s, mode: %s, node: %s\n", validators.Configuration.Network.Name, validators.Configuration.Network.Mode, validators.Configuration.Network.Node)

	index := 0
	nonce := currentNonce()
	gasPrice := validators.Staking.Gas.Price

	if validators.Configuration.Application.Infinite {
		for {
			index, nonce = PerformCreateValidators(index, nonce, gasPrice)
		}
	} else {
		pools := 1
		if validators.Configuration.Application.Count > validators.Configuration.Application.PoolSize {
			pools = int(math.RoundToEven(float64(validators.Configuration.Application.Count) / float64(validators.Configuration.Application.PoolSize)))
		}

		for poolIndex := 0; poolIndex < pools; poolIndex++ {
			if poolIndex > 1 {
				nonce := currentNonce()
				gasPrice = sdkTxs.BumpGasPrice(gasPrice)
				fmt.Println(fmt.Sprintf("Nonce is now: %d, gas price is now: %f", nonce, gasPrice))
			}

			index, nonce = PerformCreateValidators(index, nonce, gasPrice)
		}
	}
}

// PerformCreateValidators - performs the actual creation via goroutines
func PerformCreateValidators(index int, nonce int, gasPrice numeric.Dec) (int, int) {
	var waitGroup sync.WaitGroup

	for i := 0; i < validators.Configuration.Application.PoolSize; i++ {
		waitGroup.Add(1)

		go CreateValidator(index, nonce, gasPrice, &waitGroup)

		index++
		nonce++
	}

	waitGroup.Wait()

	return index, nonce
}

// CreateValidator - creates a given validator
func CreateValidator(index int, nonce int, gasPrice numeric.Dec, waitGroup *sync.WaitGroup) error {
	defer waitGroup.Done()
	accountName := fmt.Sprintf("ValidatorSpammer_Account_%d", index)
	logger.AccountLog(fmt.Sprintf("Generating a new account: %s", accountName), validators.Configuration.Application.Verbose)

	account, err := accounts.GenerateAccount(accountName)
	if err != nil {
		logger.ErrorLog(err.Error(), validators.Configuration.Application.Verbose)
		return err
	}

	logger.AccountLog(fmt.Sprintf("Generated account: %s, address: %s", account.Name, account.Address), validators.Configuration.Application.Verbose)

	fundingAccountBalance, err := balances.GetShardBalance(validators.Configuration.Funding.Account.Address, 0)
	if err != nil {
		logger.ErrorLog(fmt.Sprintf("Failed to retrieve shard balance - error: %s", err.Error()), validators.Configuration.Application.Verbose)
		return err
	}

	logger.FundingLog(fmt.Sprintf("Available funding amount in the funding account %s, address: %s is %f", validators.Configuration.Funding.Account.Name, validators.Configuration.Funding.Account.Address, fundingAccountBalance), validators.Configuration.Application.Verbose)
	fundingAmount := funding.CalculateFundingAmount(validators.Staking.Create.Validator.Amount, fundingAccountBalance, 1)
	logger.FundingLog(fmt.Sprintf("Funding account %s, address: %s with %f", account.Name, account.Address, fundingAmount), validators.Configuration.Application.Verbose)
	funding.PerformFundingTransaction(&validators.Configuration.Funding.Account, 0, account.Address, 0, fundingAmount, nonce, validators.Configuration.Funding.Gas.Limit, validators.Configuration.Funding.Gas.Price, validators.Configuration.Funding.Timeout, validators.Configuration.Funding.Attempts)

	accountStartingBalance, _ := balances.GetShardBalance(account.Address, 0)
	logger.AccountLog(fmt.Sprintf("Using account %s, address: %s to create a new validator", account.Name, account.Address), validators.Configuration.Application.Verbose)
	logger.BalanceLog(fmt.Sprintf("Account %s, address: %s has a starting balance of %f in shard %d before the test", account.Name, account.Address, accountStartingBalance, 0), validators.Configuration.Application.Verbose)

	if !accountStartingBalance.IsZero() {
		logger.TransactionLog(fmt.Sprintf("Sending create validator transaction - will wait %d seconds for it to finalize", validators.Staking.Timeout), validators.Configuration.Application.Verbose)
		fmt.Println("")

		validators.Staking.Create.Validator.Account = &account
		blsKeys := crypto.GenerateBlsKeys(validators.Staking.Create.BLSKeyCount, "")
		rawTx, err := tfStaking.CreateValidator(&account, &account, &validators.Staking, blsKeys)
		if err != nil {
			logger.ErrorLog(fmt.Sprintf("Failed to create validator - error: %s", err.Error()), validators.Configuration.Application.Verbose)
			return err
		}

		tx := sdkTxs.ToTransaction(account.Address, 0, account.Address, 0, rawTx, err)

		if validators.Staking.Timeout > 0 {
			txResultColoring := logger.ResultColoring(tx.Success, true).Render(fmt.Sprintf("tx successful: %t", tx.Success))
			logger.TransactionLog(fmt.Sprintf("Performed create validator - validator address: %s, transaction hash: %s, %s", account.Address, tx.TransactionHash, txResultColoring), validators.Configuration.Application.Verbose)
		} else {
			logger.TransactionLog(fmt.Sprintf("Performed create validator - validator address: %s, transaction hash: %s", account.Address, tx.TransactionHash), validators.Configuration.Application.Verbose)
		}

		logger.TeardownLog("Performing test teardown (returning funds and removing accounts)\n", validators.Configuration.Application.Verbose)

		if validators.Staking.Timeout > 0 {
			if tx.Success {
				exportKeys(account, blsKeys)
			}

			testing.Teardown(&account, 0, validators.Configuration.Funding.Account.Address, 0)
		} else {
			goSdkAccount.RemoveAccount(account.Name)
		}
	} else {
		logger.ErrorLog(fmt.Sprintf("Account %s, address: %s doesn't have sufficient balance to create a validator! The account should have a balance of %f but the actual balance is %f", account.Name, account.Address, fundingAmount, accountStartingBalance), validators.Configuration.Application.Verbose)
	}

	return nil
}

func currentNonce() (nonce int) {
	nonce = -1
	rpcClient, _ := validators.Configuration.Network.API.RPCClient(0)

	if validators.Configuration.Network.Mode == "local" {
		rpcNode := sdkNetwork.GenerateNodeAddress(validators.Configuration.Network.Name, "api", 0)
		remoteRPCClient, _ := sdkNetwork.NewRPCClient(rpcNode, 0)
		rpcNonce := sdkNetwork.CurrentNonce(remoteRPCClient, validators.Configuration.Funding.Account.Address)
		localNonce := sdkNetwork.CurrentNonce(rpcClient, validators.Configuration.Funding.Account.Address)

		fmt.Printf("Current RPC nonce is: %d, current local nonce is: %d\n", rpcNonce, localNonce)

		if rpcNonce > localNonce {
			nonce = int(rpcNonce)
		} else {
			nonce = int(localNonce)
		}
	} else {
		receivedNonce := sdkNetwork.CurrentNonce(rpcClient, validators.Configuration.Funding.Account.Address)
		nonce = int(receivedNonce)
	}

	fmt.Printf("The current nonce is: %d\n", nonce)

	return nonce
}

func exportKeys(account sdkAccounts.Account, blsKeys []sdkCrypto.BLSKey) error {
	dirPath := filepath.Join(validators.Configuration.BasePath, "generated", account.Address)
	keystorePath := filepath.Join(dirPath, fmt.Sprintf("%s.key", account.Address))

	if err := utils.CreateDirectory(dirPath); err != nil {
		return err
	}

	keystoreJSON, err := account.ExportKeystore(validators.Configuration.Application.Passphrase)
	if err != nil {
		return err
	}

	if len(keystoreJSON) > 0 {
		os.Remove(keystorePath)
		ioutil.WriteFile(keystorePath, keystoreJSON, 0755)
	}

	for _, blsKey := range blsKeys {
		blsKeyPath := filepath.Join(dirPath, fmt.Sprintf("%s.key", blsKey.PublicKeyHex))
		encrypted, err := blsKey.Encrypt(validators.Configuration.Application.Passphrase)

		if err == nil && len(encrypted) > 0 {
			os.Remove(blsKeyPath)
			ioutil.WriteFile(blsKeyPath, []byte(encrypted), 0755)
		}
	}

	return nil
}
