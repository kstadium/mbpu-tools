package cmd

import (


	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mbpu-tools",
		Short: "A generator for Cobra based Applications",
		Long: `Tools for MBPU.`,
	}
)


// Execute executes the root command.
func Execute() error {
	rootCmd.AddCommand(cmdMBPU)
	rootCmd.AddCommand(cmdBCCSP)
	
	cmdMBPU.AddCommand(cmdMBPUVersion)
	cmdMBPU.AddCommand(cmdMBPUTest)

	cmdBCCSP.AddCommand(cmdBCCSPVersion)
	cmdBCCSP.AddCommand(cmdBCCSPTest)
	
	return rootCmd.Execute()
}
