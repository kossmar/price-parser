package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}

	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v, a, b")
	}
	t.Fatal(message)
}

func TestSetupOutputToFile(t *testing.T) {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := usr.HomeDir + "/Documents/price-parser_test"

	file, err := setupOutputFile(dir)
	if file == nil && err != nil {
		t.Errorf("Output file was not setup, got: %d, expected: a file??.", file)
	}

	os.Remove(dir)
}

func TestOutputToFile(t *testing.T) {
	// set up and tear down test file ??
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := usr.HomeDir + "/Documents/price-parser_test"

	file, err := setupOutputFile(dir)
	newWriter := bufio.NewWriter(file)
	testString := "If you don't know, now you know."
	err1 := outputToFile(testString, newWriter)
	if err1 != nil {
		t.Errorf("Output to file failed")
	}

	os.Remove(dir)

}
