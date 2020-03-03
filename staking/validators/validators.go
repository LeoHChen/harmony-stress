package validators

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
	"github.com/SebastianJ/harmony-stress/utils"
	"github.com/SebastianJ/harmony-tf/accounts"
	"github.com/SebastianJ/harmony-tf/balances"
	"github.com/SebastianJ/harmony-tf/crypto"
	"github.com/SebastianJ/harmony-tf/funding"
	"github.com/SebastianJ/harmony-tf/logger"
	tfStaking "github.com/SebastianJ/harmony-tf/staking"
	"github.com/SebastianJ/harmony-tf/testing"
	goSdkAccount "github.com/harmony-one/go-sdk/pkg/account"
)

// CreateValidators - mass creates validators
func CreateValidators() {
	fmt.Printf("Starting validator spammer - network: %s, mode: %s, node: %s\n", Configuration.Network.Name, Configuration.Network.Mode, Configuration.Network.Node)

	index := 0
	nonce := currentNonce()

	if Configuration.Application.Infinite {
		for {
			index, nonce = PerformCreateValidators(index, nonce)
		}
	} else {
		pools := 1
		if Configuration.Application.Count > Configuration.Application.PoolSize {
			pools = int(math.RoundToEven(float64(Configuration.Application.Count) / float64(Configuration.Application.PoolSize)))
		}

		for poolIndex := 0; poolIndex < pools; poolIndex++ {
			index, nonce = PerformCreateValidators(index, nonce)
		}
	}
}

// PerformCreateValidators - performs the actual creation via goroutines
func PerformCreateValidators(index int, nonce int) (int, int) {
	var waitGroup sync.WaitGroup

	for i := 0; i < Configuration.Application.PoolSize; i++ {
		waitGroup.Add(1)

		go CreateValidator(index, nonce, &waitGroup)

		index++
		nonce++
	}

	waitGroup.Wait()

	return index, nonce
}

// CreateValidator - creates a given validator
func CreateValidator(index int, nonce int, waitGroup *sync.WaitGroup) error {
	defer waitGroup.Done()
	accountName := fmt.Sprintf("ValidatorSpammer_Account_%d", index)
	logger.AccountLog(fmt.Sprintf("Generating a new account: %s", accountName), Configuration.Application.Verbose)

	account, err := accounts.GenerateAccount(accountName)
	if err != nil {
		logger.ErrorLog(err.Error(), Configuration.Application.Verbose)
		return err
	}

	logger.AccountLog(fmt.Sprintf("Generated account: %s, address: %s", account.Name, account.Address), Configuration.Application.Verbose)

	fundingAccountBalance, err := balances.GetShardBalance(Configuration.Funding.Account.Address, 0)
	if err != nil {
		logger.ErrorLog(fmt.Sprintf("Failed to retrieve shard balance - error: %s", err.Error()), Configuration.Application.Verbose)
		return err
	}

	logger.FundingLog(fmt.Sprintf("Available funding amount in the funding account %s, address: %s is %f", Configuration.Funding.Account.Name, Configuration.Funding.Account.Address, fundingAccountBalance), Configuration.Application.Verbose)
	fmt.Printf("\n\nStaking.Create.Validator: %+v, fundingAccountBalance: %f\n\n", Staking.Create.Validator, fundingAccountBalance)
	fundingAmount := funding.CalculateFundingAmount(Staking.Create.Validator.Amount, fundingAccountBalance, 1)
	logger.FundingLog(fmt.Sprintf("Funding account %s, address: %s with %f", account.Name, account.Address, fundingAmount), Configuration.Application.Verbose)
	funding.PerformFundingTransaction(&Configuration.Funding.Account, 0, account.Address, 0, fundingAmount, nonce, Configuration.Funding.Gas.Limit, Configuration.Funding.Gas.Price, Configuration.Funding.Timeout, Configuration.Funding.Attempts)

	accountStartingBalance, _ := balances.GetShardBalance(account.Address, 0)
	logger.AccountLog(fmt.Sprintf("Using account %s, address: %s to create a new validator", account.Name, account.Address), Configuration.Application.Verbose)
	logger.BalanceLog(fmt.Sprintf("Account %s, address: %s has a starting balance of %f in shard %d before the test", account.Name, account.Address, accountStartingBalance, 0), Configuration.Application.Verbose)

	if !accountStartingBalance.IsZero() {
		logger.TransactionLog(fmt.Sprintf("Sending create validator transaction - will wait %d seconds for it to finalize", Staking.Timeout), Configuration.Application.Verbose)
		fmt.Println("")

		Staking.Create.Validator.Account = &account
		blsKeys := crypto.GenerateBlsKeys(Staking.Create.BLSKeyCount, "")
		rawTx, err := tfStaking.CreateValidator(&account, nil, &Staking, blsKeys)
		if err != nil {
			logger.ErrorLog(fmt.Sprintf("Failed to create validator - error: %s", err.Error()), Configuration.Application.Verbose)
			return err
		}

		tx := sdkTxs.ToTransaction(account.Address, 0, account.Address, 0, rawTx, err)

		if Staking.Timeout > 0 {
			txResultColoring := logger.ResultColoring(tx.Success, true).Render(fmt.Sprintf("tx successful: %t", tx.Success))
			logger.TransactionLog(fmt.Sprintf("Performed create validator - transaction hash: %s, %s", tx.TransactionHash, txResultColoring), Configuration.Application.Verbose)
		} else {
			logger.TransactionLog(fmt.Sprintf("Performed create validator - transaction hash: %s", tx.TransactionHash), Configuration.Application.Verbose)
		}

		logger.TeardownLog("Performing test teardown (returning funds and removing accounts)\n", Configuration.Application.Verbose)

		if Staking.Timeout > 0 {
			if tx.Success {
				exportKeys(account, blsKeys)
			}

			testing.Teardown(&account, 0, Configuration.Funding.Account.Address, 0)
		} else {
			goSdkAccount.RemoveAccount(account.Name)
		}
	} else {
		logger.ErrorLog(fmt.Sprintf("Account %s, address: %s doesn't have sufficient balance to create a validator! The account should have a balance of %f but the actual balance is %f", account.Name, account.Address, fundingAmount, accountStartingBalance), Configuration.Application.Verbose)
	}

	return nil
}

func currentNonce() (nonce int) {
	nonce = -1
	rpcClient, _ := Configuration.Network.API.RPCClient(0)

	if Configuration.Network.Mode == "local" {
		rpcNode := sdkNetwork.GenerateNodeAddress(Configuration.Network.Name, "api", 0)
		remoteRPCClient, _ := sdkNetwork.NewRPCClient(rpcNode, 0)
		rpcNonce := sdkNetwork.CurrentNonce(remoteRPCClient, Configuration.Funding.Account.Address)
		localNonce := sdkNetwork.CurrentNonce(rpcClient, Configuration.Funding.Account.Address)

		fmt.Printf("Current RPC nonce is: %d, current local nonce is: %d\n", rpcNonce, localNonce)

		if rpcNonce > localNonce {
			nonce = int(rpcNonce)
		} else {
			nonce = int(localNonce)
		}
	} else {
		receivedNonce := sdkNetwork.CurrentNonce(rpcClient, Configuration.Funding.Account.Address)
		nonce = int(receivedNonce)
	}

	fmt.Printf("The current nonce is: %d\n", nonce)

	return nonce
}

func exportKeys(account sdkAccounts.Account, blsKeys []sdkCrypto.BLSKey) error {
	dirPath := filepath.Join(Configuration.BasePath, "generated", account.Address)
	keystorePath := filepath.Join(dirPath, fmt.Sprintf("%s.key", account.Address))

	if err := utils.CreateDirectory(dirPath); err != nil {
		return err
	}

	keystoreJSON, err := account.ExportKeystore(Configuration.Application.Passphrase)
	if err != nil {
		return err
	}

	if len(keystoreJSON) > 0 {
		os.Remove(keystorePath)
		ioutil.WriteFile(keystorePath, keystoreJSON, 0755)
	}

	for _, blsKey := range blsKeys {
		blsKeyPath := filepath.Join(dirPath, fmt.Sprintf("%s.key", blsKey.PublicKeyHex))
		encrypted, err := blsKey.Encrypt(Configuration.Application.Passphrase)

		if err == nil && len(encrypted) > 0 {
			os.Remove(blsKeyPath)
			ioutil.WriteFile(blsKeyPath, []byte(encrypted), 0755)
		}
	}

	return nil
}
