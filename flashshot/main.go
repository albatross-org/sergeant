package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/gen2brain/beeep"
	"github.com/sirupsen/logrus"
)

// dateFormat used for saving screenshots.
var dateFormat = "20060102-1504-05"

// tempDir creates a temporary directory and returns a function that will remove the temporary directory.
// Instead of using ioutil.TempDir, we generate one ourselves since we need it to have lots of permissions.
func tempDir() (path string, cleanup func()) {
	path = filepath.Join(os.TempDir(), fmt.Sprintf("flashshot%d", rand.Intn(9999999999)))

	err := os.Mkdir(path, 0755)
	if err != nil {
		panic(fmt.Errorf("could not create temporary directory: %s", err))
	}

	err = os.Mkdir(filepath.Join(path, "questions"), 0755)
	if err != nil {
		panic(fmt.Errorf("couldn't create temporary 'questions' directory: %s", err))
	}

	err = os.Mkdir(filepath.Join(path, "answers"), 0755)
	if err != nil {
		panic(fmt.Errorf("couldn't create temporary 'answers' directory: %s", err))
	}

	return path, func() {
		err = os.RemoveAll(path)
		if err != nil {
			panic(fmt.Errorf("could not remove temporary directory: %s", err))
		}
	}
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

// cleanScreenshotDir will remove any temporary screenshot files from a directory. Ut does this by removing and then recreating
// the directory. This isn't ideal.
func cleanScreenshotDir(dir string) error {
	err := os.RemoveAll(filepath.Join(dir, "questions"))
	if err != nil {
		return fmt.Errorf("couldn't remove temporary 'questions' directory: %s", err)
	}

	err = os.RemoveAll(filepath.Join(dir, "answers"))
	if err != nil {
		return fmt.Errorf("couldn't remove temporary 'answers' directory: %s", err)
	}

	err = os.Mkdir(filepath.Join(dir, "questions"), 0755)
	if err != nil {
		return fmt.Errorf("couldn't recreate temporary 'questions' directory: %s", err)
	}

	err = os.Mkdir(filepath.Join(dir, "answers"), 0755)
	if err != nil {
		return fmt.Errorf("couldn't recreate temporary 'answers' directory: %s", err)
	}

	return nil
}

// screenshot takes a screenshot using the maim tool.
func screenshot(dest string) error {
	maimCmd := exec.Command("maim", "-s", dest)

	bytes, err := maimCmd.CombinedOutput()
	output := string(bytes)

	if err != nil && !strings.Contains(output, "Selection was cancelled by keystroke or right-click.") {
		return fmt.Errorf("maim command '%s' exited with message '%s', error: %w", maimCmd.String(), output, err)
	}

	return nil
}

// tileshot creates a image from multiple screenshots put together.
// It will keep stiching selections together until it recieves a value on the "done" channel.
func tileshot(destination, dir string, done chan bool) {
	names := []string{}

	for {
		select {
		case <-done:
			err := tileImages(destination, names)
			if err != nil {
				logrus.Error("Couldn't tile screenshot together:")
				logrus.Fatal(err)
			}

			err = cleanScreenshotDir(dir)
			if err != nil {
				logrus.Error("Error cleaning up temporary screenshot dir:")
				logrus.Fatal(err)
			}

			return
		default:
			filename := fmt.Sprintf("flashshot-%s.png", time.Now().Format(dateFormat))
			path := filepath.Join(dir, "questions", filename)
			err := screenshot(path)

			if err != nil {
				logrus.Error("Error taking screenshot:")
				logrus.Fatal(err)
			}
		}
	}
}

func getEvents() (chan keylogger.InputEvent, func() error, error) {
	// find keyboard device, does not require a root permission
	keyboard := keylogger.FindKeyboardDevice()

	// check if we found a path to keyboard
	if len(keyboard) <= 0 {
		return nil, nil, fmt.Errorf("no keyboard found...you will need to provide manual input path")
	}

	logrus.Println("Found a keyboard at", keyboard)

	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		return nil, nil, err
	}

	return k.Read(), k.Close, err
}

func main() {
	rand.Seed(time.Now().Unix())

	dir, _ := tempDir()
	// defer cleanup()

	events, close, err := getEvents()
	defer close()
	if err != nil {
		logrus.Error("Couldn't get the keyboard events:")
		logrus.Fatal(err)
	}

	doneChan := make(chan bool)
	state := "idle"

	// TODO: you're overcomplicating things by making tileshot a seperate function. It would make much more sense if the logic was here.

	// range of events
	for e := range events {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:
			if e.KeyPress() {
				switch e.KeyString() {
				case "Q":
					if state == "question" {
						logrus.Println("Can't start screenshotting question, that's the current state.")
						continue
					}

					state = "question"

					beeep.Notify("flashshot", "Scan in your question! Press A when you're ready to scan in the answer.", "")

					go tileshot("question.png", dir, doneChan)
				case "A":
					logrus.Println("A pressed!")

					if state != "question" {
						logrus.Println("Can't start scanning in an answer if there's no question for it to follow.")
						continue
					}

					state = "answer"

					doneChan <- true
					beeep.Notify("flashshot", "Scan in your answer! Press D when you're done.", "")
					go tileshot("answer.png", dir, doneChan)

				case "D":
					logrus.Println("D pressed!")
					state = "done"
					doneChan <- true
				}
			}
		}
	}
}
