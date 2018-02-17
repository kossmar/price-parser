package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verboseCmd = &cobra.Command{
	Use:   "verbose",
	Short: "Show more price information",
	Long:  `I don't know what else to tell you`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("so verbose")
	},
}

func init() {
	ParsePriceCmd.AddCommand(verboseCmd)
}
