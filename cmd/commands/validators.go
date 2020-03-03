package commands

import (
	"path/filepath"

	cmdConfig "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/SebastianJ/harmony-stress/staking/validators"
	"github.com/spf13/cobra"
)

func init() {
	cmdValidators := &cobra.Command{
		Use:   "validators",
		Short: "Stress test validator creation",
		Long:  "Stress tests the creation of validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stressTestValidators(cmd)
		},
	}

	RootCmd.AddCommand(cmdValidators)
}

func stressTestValidators(cmd *cobra.Command) error {
	basePath, err := filepath.Abs(cmdConfig.Persistent.Path)
	if err != nil {
		return err
	}

	if err := validators.Configure(basePath, cmdConfig.Persistent); err != nil {
		return err
	}

	validators.CreateValidators()

	return nil
}
