package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/sergeant"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:       "query [views, paths, sets, config]",
	Short:     "List the views, paths or sets currently loaded in the program",
	Long:      ``,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"views", "paths", "sets", "config"},

	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}

		config, err := sergeant.LoadConfig(configPath)
		if err != nil {
			logrus.Fatal(err)
		}

		underlyingStore, err := albatross.FromConfig(config.Store)
		if err != nil {
			logrus.Fatal(err)
		}

		store := sergeant.NewStore(underlyingStore, config)
		if err != nil {
			logrus.Fatal(err)
		}

		switch args[0] {
		case "views":
			for viewName := range sergeant.DefaultViews {
				fmt.Println(viewName)
			}
		case "paths":
			all, _, err := store.Set("all")
			if err != nil {
				logrus.Fatal(err)
			}

			for _, card := range all.Cards {
				fmt.Println(card.Path)
			}
		case "sets":
			for set := range store.Sets {
				fmt.Println(set)
			}
		case "config":
			bytes, err := json.Marshal(store.Config)
			if err != nil {
				logrus.Fatal(err)
			}

			os.Stdout.Write(bytes)
		}

	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
