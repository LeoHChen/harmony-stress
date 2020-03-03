package commands

import (
	"fmt"
	"os"

	cmd "github.com/SebastianJ/harmony-stress/config/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// VersionWrap - version displayed in case of errors
	VersionWrap = ""

	// RootCmd - main entry point for Cobra commands
	RootCmd = &cobra.Command{
		Use:          "stress",
		Short:        "Stress test",
		SilenceUsage: true,
		Long:         "Harmony stress test tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}
)

func init() {
	cmd.Persistent = cmd.PersistentFlags{}
	RootCmd.PersistentFlags().StringVar(&cmd.Persistent.Path, "path", ".", "<path>")
	RootCmd.PersistentFlags().StringVar(&cmd.Persistent.Network, "network", "localnet", "<network>")
	RootCmd.PersistentFlags().StringVar(&cmd.Persistent.NetworkMode, "network-mode", "api", "<mode>")
	RootCmd.PersistentFlags().StringVar(&cmd.Persistent.Node, "node", "", "<node>")
	RootCmd.PersistentFlags().StringVar(&cmd.Persistent.ApplicationMode, "app-mode", "async", "<app-mode>")
	RootCmd.PersistentFlags().StringVar(&cmd.Persistent.From, "from", "", "<from>")
	RootCmd.PersistentFlags().StringVar(&cmd.Persistent.Passphrase, "passphrase", "", "<passphrase>")
	RootCmd.PersistentFlags().BoolVar(&cmd.Persistent.Infinite, "infinite", false, "<infinite>")
	RootCmd.PersistentFlags().IntVar(&cmd.Persistent.Count, "count", 1000, "<count>")
	RootCmd.PersistentFlags().IntVar(&cmd.Persistent.PoolSize, "pool-size", 100, "<pool-size>")
	RootCmd.PersistentFlags().IntVar(&cmd.Persistent.Timeout, "timeout", 0, "<pool-size>")
	RootCmd.PersistentFlags().BoolVar(&cmd.Persistent.Verbose, "verbose", false, "<verbose>")
	RootCmd.PersistentFlags().BoolVar(&cmd.Persistent.VerboseGoSdk, "verbose-go-sdk", false, "<verbose-go-sdk>")
	RootCmd.PersistentFlags().IntVar(&cmd.Persistent.PprofPort, "pprof-port", -1, "<pprof-port>")
}

// Execute kicks off the hmy CLI
func Execute() {
	RootCmd.SilenceErrors = true
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(errors.Wrapf(err, "commit: %s, error", VersionWrap).Error())
		os.Exit(1)
	}
}
