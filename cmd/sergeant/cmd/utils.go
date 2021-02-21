package cmd

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// letterBytes are the letters used to generate a random string.
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// randomString generates a string consisting of characters from letterBytes that is n characters long.
// Courtesy: https://stackoverflow.com/a/31832326
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// exists reports whether the named file or directory exists.
func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// checkFlag checks an error from a call to cmd.Flags().GetXXX and displays an error if there is one.
func checkFlag(err error, flag, cmd string) {
	if err != nil {
		fmt.Printf("Error getting flag %s from command %s:\n", flag, cmd)
		fmt.Println(err)
		os.Exit(1)
	}
}

// tempDir creates a temporary directory and returns a function that will remove the temporary directory.
// Instead of using ioutil.TempDir, we generate one ourselves since we need it to have lots of permissions.
func tempDir() (path string, cleanup func()) {
	path, err := ioutil.TempDir("", "sergeant-")
	if err != nil {
		panic(fmt.Errorf("could not create temporary directory: %s", err))
	}

	return path, func() {
		err = os.RemoveAll(path)
		if err != nil {
			panic(fmt.Errorf("could not remove temporary directory: %s", err))
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
