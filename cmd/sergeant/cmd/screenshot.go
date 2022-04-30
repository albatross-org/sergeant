package cmd

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/sergeant"
	"github.com/fatih/color"
	hook "github.com/robotn/gohook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// screenshotCmd represents the 'screenshot' command.
var screenshotCmd = &cobra.Command{
	Use:   "screenshot --path [path]",
	Short: "Quickly scan cards into the program",
	Long: `Screenshot lets you quickly add cards to the program by screenshotting questions and answers.
	
For example:

	$ sergeant screenshot --path 'further-maths/core-pure-1/chapter-1-complex-numbers/ex1a'
	# Listens for keyboard "Q" (question), "A" (answer), "D" (done) and "C" (cancel).
	
This command will add different cards to the path 'further-maths/core-pure-1/chapter-1-complex-numbers/ex1a'.

You could think of this command as automating the process of:

	- Taking screenshots manually
	- Stiching together the ones that run over page breaks
	- Running the 'sergeant add' command.
	
It allows you to turn minutes per question into seconds.

You can also specify tags:

	$ sergeant screenshot \ 
		--path 'further-maths/core-pure-1/chapter-1-complex-numbers/ex1a' \
		--tags "@?school" --tags "@?further-maths"
		
Which will be present in every entry.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		checkFlag(err, "--path", "screenshot")

		tags, err := cmd.Flags().GetStringSlice("tags")
		checkFlag(err, "--tags", "screenshot")

		disableGit, err := cmd.Flags().GetBool("disable-git")
		checkFlag(err, "--disable-git", "screenshot")

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

		if disableGit {
			store.DisableGit()
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/profile", pprof.Profile)
		go http.ListenAndServe(":7777", mux)

		tempPath, cleanup := tempDir()
		images := cardImages{tempPath: tempPath, cleanup: cleanup, mu: &sync.Mutex{}}

		bold := color.New(color.Bold)
		italic := color.New(color.Italic)
		yellow := color.New(color.FgHiYellow)

		fmt.Println("Starting screenshotting program:")
		fmt.Println("")

		fmt.Print("- ", bold.Sprint("SHIFT+Q"), ": Scan a question.\n")
		fmt.Print("- ", bold.Sprint("SHIFT+A"), ": Scan an answer.\n")
		fmt.Print("- ", bold.Sprint("SHIFT+D"), ": Finish the current card and start again.\n")
		fmt.Print("- ", bold.Sprint("SHIFT+C"), ": Cancel the current card and start again.\n")
		fmt.Print("- ", bold.Sprint("SHIFT+CTRL+C"), ": Quit the program.\n")

		fmt.Println("\nTemporary work will be saved to", tempPath, "and removed afterwards.")
		fmt.Println("")

		count := 1

		hook.Register(hook.KeyDown, []string{"q", "shift"}, func(e hook.Event) {
			if len(images.answerImages) != 0 {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Creating card\n")

				questionImage, answerImage, err := images.Build()
				if err != nil {
					logrus.Fatal(err)
				}

				entryPath, err := fastCreateCard(store, path, tags, questionImage, answerImage)
				if err != nil {
					logrus.Fatal(err)
				}

				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Created: ", italic.Sprint(entryPath), "\n\n")
				count++

				err = images.Cleanup()
				if err != nil {
					logrus.Fatal(err)
				}
			}

			if len(images.questionImages) == 0 {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Scanning question\n")
			} else {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Scanning question, part ", len(images.questionImages)+1, "\n")
			}

			err = images.ScanQuestion()
			if err != nil {
				logrus.Fatal(err)
			}

			if len(images.questionImages) == 1 {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Scanned question\n")
			} else {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Scanned question, part ", len(images.questionImages), "\n")
			}
		})

		hook.Register(hook.KeyDown, []string{"a", "shift"}, func(e hook.Event) {
			if len(images.answerImages) == 0 {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Scanning answer\n")
			} else {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Scanning answer, part ", len(images.answerImages)+1, "\n")
			}

			err = images.ScanAnswer()
			if err != nil {
				logrus.Fatal(err)
			}

			if len(images.answerImages) == 1 {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Scanned answer\n")
			} else {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Scanned answer, part ", len(images.answerImages), "\n")
			}
		})

		hook.Register(hook.KeyDown, []string{"d", "shift"}, func(e hook.Event) {
			fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Creating card\n")

			questionImage, answerImage, err := images.Build()
			if err != nil {
				logrus.Fatal(err)
			}

			entryPath, err := fastCreateCard(store, path, tags, questionImage, answerImage)
			if err != nil {
				logrus.Fatal(err)
			}

			fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Created: ", italic.Sprint(entryPath), "\n\n")
			count++

			err = images.Cleanup()
			if err != nil {
				logrus.Fatal(err)
			}
		})

		hook.Register(hook.KeyDown, []string{"c", "shift"}, func(e hook.Event) {
			fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Cancelling\n")

			err = images.Cancel()
			if err != nil {
				logrus.Fatal(err)
			}
		})

		hook.Register(hook.KeyDown, []string{"c", "ctrl", "shift"}, func(e hook.Event) {
			fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Quiting program\n")
			hook.End()

			// Save the last flashcard created.
			if len(images.answerImages)+len(images.questionImages) != 0 {
				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Creating last card\n")
				questionImage, answerImage, err := images.Build()
				if err != nil {
					logrus.Fatal(err)
				}

				entryPath, err := createCard(store, path, tags, questionImage, answerImage)
				if err != nil {
					logrus.Fatal(err)
				}

				fmt.Print(bold.Sprint("Card("), yellow.Sprint(count), bold.Sprint(")"), " - Created: ", italic.Sprint(entryPath), "\n")
			}

			cleanup()
		})

		s := hook.Start()
		<-hook.Process(s)
	},
}

// cardImages manages creating a question and answer image from multiple screenshots.
type cardImages struct {
	tempPath string
	cleanup  func()

	mu *sync.Mutex

	questionImages []string
	answerImages   []string

	finalQuestionImage string
	finalAnswerImage   string
}

// path returns the path relative to the temporary dir.
func (c *cardImages) path(path string) string {
	return filepath.Join(c.tempPath, path)
}

// generateFinalQuestionImage sets the finalQuestionImage by creating appending the images in questionImages together.
func (c *cardImages) generateFinalQuestionImage() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.questionImages) == 0 {
		return fmt.Errorf("there are no questionImages to tile together")
	}

	dest := c.path("question.png")
	err := tileImages(dest, c.questionImages)
	if err != nil {
		return err
	}

	c.finalQuestionImage = dest
	return nil
}

// generateFinalAnswerImage sets the finalAnswerImage by creating appending the images in answerImages together.
func (c *cardImages) generateFinalAnswerImage() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.answerImages) == 0 {
		return fmt.Errorf("there are no answerImages to tile together")
	}

	dest := c.path("answer.png")
	err := tileImages(dest, c.answerImages)
	if err != nil {
		return err
	}

	c.finalAnswerImage = dest
	return nil
}

// ScanQuestion scans a question in.
func (c *cardImages) ScanQuestion() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dest := c.path("question-" + randomString(16) + ".png")

	err := screenshot(dest)
	if err != nil {
		return err
	}

	c.questionImages = append(c.questionImages, dest)
	return nil
}

// ScanAnswer scans an answer in.
func (c *cardImages) ScanAnswer() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dest := c.path("answer-" + randomString(16) + ".png")

	err := screenshot(dest)
	if err != nil {
		return err
	}

	c.answerImages = append(c.answerImages, dest)
	return nil
}

// Cancel cancels the current card scan.
func (c *cardImages) Cancel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, intermediateQuestionImage := range c.questionImages {
		err := os.Remove(intermediateQuestionImage)
		if err != nil {
			return err
		}
	}

	for _, intermediateAnswerImage := range c.answerImages {
		err := os.Remove(intermediateAnswerImage)
		if err != nil {
			return err
		}
	}

	c.questionImages = []string{}
	c.answerImages = []string{}

	return nil
}

// Build returns the final question image and the final answer image and resets all values internally.
func (c *cardImages) Build() (questionImage string, answerImage string, err error) {
	err = c.generateFinalQuestionImage()
	if err != nil {
		return "", "", err
	}

	err = c.generateFinalAnswerImage()
	if err != nil {
		return "", "", err
	}

	questionImage = c.finalQuestionImage
	answerImage = c.finalAnswerImage

	return questionImage, answerImage, nil
}

// Cleanup deletes all intermediate workings. This should be called after a successful call to .Build().
func (c *cardImages) Cleanup() error {
	err := c.Cancel()
	if err != nil {
		return err
	}

	err = os.Remove(c.finalQuestionImage)
	if err != nil {
		return err
	}

	err = os.Remove(c.finalAnswerImage)
	if err != nil {
		return err
	}

	return nil
}

// tileImages takes a list of paths and outputs an image with all those paths tiled together. It relies on having
// ImageMagic installed.
func tileImages(destination string, input []string) error {
	args := append(input, "-append", destination)
	convertCmd := exec.Command("convert", args...)

	bs, err := convertCmd.Output()
	if err != nil {
		return fmt.Errorf("convert command '%s' destination '%s', an error: %w", convertCmd.String(), string(bs), err)
	}

	return nil
}

func init() {
	screenshotCmd.Flags().StringP("path", "p", "", "path to where the card should go")
	screenshotCmd.Flags().StringSliceP("tags", "t", []string{}, "tags to add to the entry created")
	screenshotCmd.Flags().Bool("disable-git", false, "don't create commits for every question scanned")

	rootCmd.AddCommand(screenshotCmd)
}
