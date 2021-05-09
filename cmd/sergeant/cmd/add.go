package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/sergeant"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// addCmd represents the 'add' command.
var addCmd = &cobra.Command{
	Use:   "add --path [path] --question [question image] --answer [answer image]",
	Short: "Add a card",
	Long: `Add lets you add a card to the database of questions.
	
For example:

	$ sergeant add --path 'further-maths/core-pure-1/chapter-1-complex-numbers' --question 'question.png' --answer 'answer.png'
	# Or, using the short versions of the flags:
	$ sergeant add -p 'further-maths/core-pure-1/chapter-1-complex-numbers' -q 'question.png' -a 'answer.png'

You can also add tags:

	$ sergeant add \
		-p 'further-maths/core-pure-1/chapter-1-complex-numbers' \
		-q 'question.png' \
		-a 'answer.png' \
		-t @?school -t @?further-maths 

This is pretty longwinded and slow to use manually. If you want to scan in lots of questions very quickly, it's much easier to use
the 'screenshot' command:

	$ sergeant screenshot --path 'further-maths/core-pure-1/chapter-1-complex-numbers/ex1a'
	# Listens for keyboard "Q" (question), "A" (answer), "D" (done) and "C" (cancel)

For more information, see:

	$ sergeant screenshot --help
	`,

	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}

		config, err := sergeant.LoadConfig(configPath)
		if err != nil {
			logrus.Fatal(err)
		}

		store, err := albatross.FromConfig(config.Store)
		if err != nil {
			logrus.Fatal(err)
		}

		path, err := cmd.Flags().GetString("path")
		checkFlag(err, "--path", "add")

		questionPath, err := cmd.Flags().GetString("question")
		checkFlag(err, "--question", "add")

		answerPath, err := cmd.Flags().GetString("answer")
		checkFlag(err, "--answer", "add")

		tags, err := cmd.Flags().GetStringSlice("tags")
		checkFlag(err, "--tags", "add")

		entryPath, err := createCard(store, path, tags, questionPath, answerPath)
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Print("Success! Your new card exists at: ")
		color.New(color.Bold).Print(entryPath)
		fmt.Println("")
	},
}

// createCard creates a card with at the given path and tags with the question and answer as an attachment.
// It returns the path to the new card and an error if there was one.
func createCard(store *albatross.Store, path string, tags []string, questionPath, answerPath string) (string, error) {
	if !exists(questionPath) {
		return "", fmt.Errorf("path to question image does not exist")
	}

	if !exists(answerPath) {
		return "", fmt.Errorf("path to answer image does not exist")
	}

	// We can omit lots of fields here since they won't be used to generate the entry content.
	card := &sergeant.Card{
		ID:   randomString(16),
		Date: time.Now(),
		Tags: tags,
	}

	content, err := card.Content()
	if err != nil {
		return "", fmt.Errorf("couldn't get card content: %s", err)
	}

	// path is only something like "further-maths/core-pure-1/chapter-1-complex-numbers".
	// We need to create a path that's unique to the card, such as "further-maths/core-pure-1/chapter-1-complex-numbers/question-0NiDQqGdzxTSipJa".
	entryPath := filepath.Join(path, "question-"+card.ID)

	err = store.Create(entryPath, content)
	if err != nil {
		return "", fmt.Errorf("couldn't create card entry: %s", err)
	}

	// It might be that the question or answer image is called something like "screenshot.png". Entries are expected to have
	// question/answer attachments in the form "question.png" or "answer.jpg" so we need to create those here.
	attachmentQuestionPath := "question" + filepath.Ext(questionPath)
	attachmentAnswerPath := "answer" + filepath.Ext(answerPath)

	err = store.AttachCopyWithName(entryPath, questionPath, attachmentQuestionPath)
	if err != nil {
		return "", fmt.Errorf("couldn't attach question image %q: %w", questionPath, err)
	}

	err = store.AttachCopyWithName(entryPath, answerPath, attachmentAnswerPath)
	if err != nil {
		return "", fmt.Errorf("couldn't attach answer image %q: %w", answerPath, err)
	}

	return entryPath, nil
}

func init() {
	addCmd.Flags().StringP("path", "p", "", "path to where the card should go")
	addCmd.Flags().StringP("question", "q", "", "path to the question image")
	addCmd.Flags().StringP("answer", "a", "", "path to the answer image")

	addCmd.Flags().StringSliceP("tags", "t", []string{}, "tags to add to the entry created")

	rootCmd.AddCommand(addCmd)
}
