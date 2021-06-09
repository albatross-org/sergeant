package cmd

import (
	"fmt"
	"os"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/sergeant"
	"github.com/albatross-org/sergeant/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sergeant",
	Short: "sergeant is a knowledge application tool.",
	Long:  `Sergeant is a tool for practicing the application of knowledge, as a supplement to existing tools for practicing recall such as Anki or Mnemosyne.`,

	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}

		config, err := sergeant.LoadConfig(path)
		if err != nil {
			logrus.Fatal(err)
		}

		underlyingStore, err := albatross.FromConfig(config.Store)
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Println(underlyingStore.UsingGit())
		store := sergeant.NewStore(underlyingStore, config)

		set, warnings, err := store.Set("all")
		if err != nil {
			logrus.Fatal(err)
		}

		for path, warning := range warnings {
			logrus.Warningf("Malformed card: %s -> %s", path, warning)
		}

		completed := 0
		for _, card := range set.Cards {
			if card.TotalCompletions() > 0 {
				completed++
			}
		}

		logrus.Infof("Loaded %d cards, %d completd. (%.2f%%)", len(set.Cards), completed, 100*float64(completed)/float64(len(set.Cards)))

		server.Run(store)
	},
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "sergeant config, defaults to ~/.config/sergeant/config.yaml")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
