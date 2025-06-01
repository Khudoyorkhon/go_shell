package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	currentDir, dirErr := os.Getwd()
	currentUserPath, userErr := os.UserHomeDir()
	currentUser := filepath.Base(currentUserPath)

	if userErr != nil {
		fmt.Fprintln(os.Stderr, "Error getting user home directory:", userErr)
		return
	}
	if dirErr != nil {
		fmt.Fprintln(os.Stderr, "Error getting current directory:", dirErr)
		return
	}

	for {
		fmt.Printf("Current directory of user %s is: %s > ", currentUser, currentDir)
		// Read the keyboad input.
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

}

func execInput(input string) error {
	// Remove the newline character
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSpace(input)

	args := strings.Split(input, " ")
	var cmd *exec.Cmd

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return errors.New("cd command requires an argument")
		}
		return os.Chdir(args[1])
	case "hostname":
		if len(args) > 1 {
			return errors.New("hostname command does not take any arguments")
		}
		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("error getting hostname: %w", err)
		}
		fmt.Println(hostname)
		return nil
	case "exit":
		os.Exit(0)
	}

	// Prepare the command to execute
	if runtime.GOOS == "windows" {
		full := append([]string{"/C", args[0]}, args[1:]...)
		cmd = exec.Command("cmd", full...)
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}

	// Set the correct output device
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command
	return cmd.Run()
}
