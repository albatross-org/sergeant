package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/sergeant"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// completeCmd represents the 'complete' command.
var completeCmd = &cobra.Command{
	Use:   "complete --path [path to question] --time [amount of time taken] (perfect|minor|major)",
	Short: "Complete a card",
	Long: `Complete lets you manually add a completion to a card using the command line rather than the Web UI.
	
For example:

	$ sergeant complete perfect --path 'further-maths/core-pure-1/chapter-1-complex-numbers/mixed-exercise-1/question-abcdef' --time "3m47s" 
	# Or, using the short versions of the flags:
	$ sergeant complete perfect -p 'further-maths/core-pure-1/chapter-1-complex-numbers/mixed-exercise-1/question-abcdef' -t "3m47s"
	`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"perfect", "minor", "major"},

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

		path, err := cmd.Flags().GetString("path")
		checkFlag(err, "--path", "complete")

		timeTaken, err := cmd.Flags().GetDuration("time")
		checkFlag(err, "--time", "complete")

		err = store.AddCompletion(path, args[0], sergeant.Completion{
			Date:     time.Now(),
			Duration: timeTaken,
		})
		if err != nil {
			fmt.Printf("Error adding %q completion to card %q: %s\n", args[0], path, err)
			os.Exit(1)
		}

		fmt.Printf("Success! Added %q completion in %s to card:\n", args[0], timeTaken)
		color.New(color.Bold).Print(path)
		fmt.Println("")
	},
}

func init() {
	completeCmd.Flags().StringP("path", "p", "", "path to the card")
	completeCmd.Flags().DurationP("time", "t", time.Duration(0), "time taken to complete the card, in XhYmZs or YmZs format")

	rootCmd.AddCommand(completeCmd)
}
