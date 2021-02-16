package cmd

import (
	"fmt"
	"os"

	"github.com/albatross-org/sergeant/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sergeant",
	Short: "sergeant is a knowledge application tool.",
	Long:  `Sergeant is a tool for practicing the application of knowledge, as a supplement to existing tools for practicing recall such as Anki or Mnemosyne.`,

	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize()
}
