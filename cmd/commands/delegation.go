package commands

import (
	"path/filepath"

	cmdConfig "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/SebastianJ/harmony-stress/staking/delegation"
	"github.com/spf13/cobra"
)

func init() {
	cmdDelegation := &cobra.Command{
		Use:   "delegations",
		Short: "Stress test delegation transactions",
		Long:  "Stress test delegation transactions (delegate / undelegate)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stressTestDelegations(cmd)
		},
	}

	cmdConfig.Delegation = cmdConfig.DelegationFlags{}
	cmdDelegation.Flags().StringVar(&cmdConfig.Delegation.ValidatorAddress, "validator-address", "", "<amount>")
	cmdDelegation.Flags().IntVar(&cmdConfig.Delegation.FromShardID, "from-shard", 0, "<shardID>")
	cmdDelegation.Flags().IntVar(&cmdConfig.Delegation.ToShardID, "to-shard", 0, "<shardID>")
	cmdDelegation.Flags().StringVar(&cmdConfig.Delegation.Amount, "amount", "", "<amount>")
	cmdDelegation.Flags().BoolVar(&cmdConfig.Delegation.OnlyActive, "active", true, "<active>")

	RootCmd.AddCommand(cmdDelegation)
}

func stressTestDelegations(cmd *cobra.Command) error {
	basePath, err := filepath.Abs(cmdConfig.Persistent.Path)
	if err != nil {
		return err
	}

	if err := delegation.Configure(basePath, cmdConfig.Persistent, cmdConfig.Delegation); err != nil {
		return err
	}

	delegation.StressTestDelegations(cmdConfig.Delegation.ValidatorAddress, cmdConfig.Delegation.OnlyActive)

	return nil
}
