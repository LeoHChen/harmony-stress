package commands

import (
	"path/filepath"

	cmdConfig "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/SebastianJ/harmony-stress/transactions"
	"github.com/spf13/cobra"
)

func init() {
	cmdTxs := &cobra.Command{
		Use:   "txs",
		Short: "Stress test normal transactions",
		Long:  "Stress test normal transactions with or without tx payloads",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stressTestTransactions(cmd)
		},
	}

	cmdConfig.Tx = cmdConfig.TxFlags{}
	cmdTxs.Flags().StringVar(&cmdConfig.Tx.Mode, "mode", "async", "<tx-mode>")
	cmdTxs.Flags().IntVar(&cmdConfig.Tx.FromShardID, "from-shard", 0, "<shardID>")
	cmdTxs.Flags().IntVar(&cmdConfig.Tx.ToShardID, "to-shard", 0, "<shardID>")
	cmdTxs.Flags().StringVar(&cmdConfig.Tx.Amount, "amount", "", "<amount>")
	cmdTxs.Flags().StringVar(&cmdConfig.Tx.GasPrice, "gas-price", "", "<gas-price>")

	RootCmd.AddCommand(cmdTxs)
}

func stressTestTransactions(cmd *cobra.Command) error {
	basePath, err := filepath.Abs(cmdConfig.Persistent.Path)
	if err != nil {
		return err
	}

	if err := transactions.Configure(basePath, cmdConfig.Persistent, cmdConfig.Tx); err != nil {
		return err
	}

	if transactions.Configuration.Transactions.Mode == "sync" {
		_, err := transactions.BulkSendTransactions()
		if err != nil {
			return err
		}
	} else if transactions.Configuration.Transactions.Mode == "async" {
		transactions.AsyncBulkSendTransactions()
	}

	return nil
}
