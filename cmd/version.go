package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Long:  `Print version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(C.Red("0.0.0"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
