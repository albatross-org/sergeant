package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/sergeant"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// devCmd represents the 'dev' command.
var devCmd = &cobra.Command{
	Use:   "dev [subcommand]",
	Short: "Test and debug the program",
	Long:  `Dev contains small programs that are currently being worked on.`,

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

		set, warnings, err := store.Set("all")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for path, warning := range warnings {
			fmt.Println("Malformed card:", path, "->", warning)
		}

		difficulties := sergeant.NewViewDifficulties(time.Now().UnixNano())
		pathTrie, _ := difficulties.BuildTrie(set)

		// fmt.Println("digraph {")
		// fmt.Println("\trankdir=LR")
		// pathTrie.Walk(func(child string, value interface{}) error {
		// 	parent := pathTrie.Get(filepath.Dir(child))
		// 	if parent == nil {
		// 		return nil
		// 	}

		// 	diff := value.(*sergeant.ProbabilityNode).Difficulty
		// 	fmt.Printf("\t\"%s\" -> \"%s\" [label=\"%f\" fontsize=%f]\n", parent.(*sergeant.ProbabilityNode).Path, child, diff, math.Pow((1-diff), 2)*64)

		// 	return nil
		// })
		// fmt.Println("}")

		fmt.Println("G = nx.DiGraph()")
		pathTrie.Walk(func(child string, value interface{}) error {
			fmt.Printf("G.add_node('%s', questions_completed=0)\n", child)
			fmt.Printf("G.add_edge('%s', '%s')\n", filepath.Dir(child), child)

			// fmt.Printf("\t\"%s\" -> \"%s\" [label=\"%f\" fontsize=%f]\n", parent.(*sergeant.ProbabilityNode).Path, child, diff, math.Pow((1-diff), 2)*64)

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
