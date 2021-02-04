package main

import (
	"fmt"

	"github.com/albatross-org/pelican"
	"github.com/sirupsen/logrus"
)

func main() {
	set, err := pelican.Load("pelican", "")
	if err != nil {
		logrus.Fatal(err)
	}

	// out, err := json.MarshalIndent(set, "", " ")
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// _ = out
	// fmt.Println(string(out))

	flashcards := set.List()
	for _, flashcard := range flashcards {
		fmt.Println(flashcard.Name())
	}
}
