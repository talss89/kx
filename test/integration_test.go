package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	err := os.Chdir("../bin")

	if err != nil {
		fmt.Println("Failed to change directory:", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()

	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		os.Exit(1)
	}

	BINPATH = dir + "/kx-cov"

	code := m.Run()
	// Teardown code here
	os.Exit(code)
}

func TestCliHelp(t *testing.T) {
	output, err := runBinary([]string{"help"})
	if err != nil {
		t.Fatal("Failed to run binary:", err)
	}

	if !bytes.Contains(output, []byte("USAGE:")) {
		t.Errorf("Expected usage information, got:\n%s", output)
	}
}

func TestShellTimeout(t *testing.T) {
	output, err := runBinary([]string{"1s", "--shell", "integration", "--ctx", "example-context-2"})
	fmt.Print(string(output))
	if err != nil {
		t.Fatal("Failed to run binary:", err)
	}

	if output == nil || !bytes.Contains(output, []byte("Session expired")) {
		t.Errorf("Expected session expiration message, got:\n%s", output)
	}

}
