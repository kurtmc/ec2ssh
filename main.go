package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/kurtmc/ec2search/search"
	"golang.org/x/crypto/ssh/terminal"
)

var mutex = &sync.Mutex{}

type PrintJob struct {
	Host    string
	Message string
}

func main() {
	instances, err := search.ListInstances(os.Args[1])
	n := len(instances)
	if err != nil {
		panic(err)
	}
	in := make(chan string, n)
	out := make(chan PrintJob, n)
	for _, i := range instances {
		in <- i
	}

	for i := 0; i < 10; i++ {
		go func() {
			host := <-in
			commandOutput, err := RunCommand(host, os.Args[2])
			if err != nil {
				panic(err)
			}
			out <- PrintJob{Host: host, Message: commandOutput}
		}()
	}

	for i := 0; i < n; i++ {
		printJob := <-out
		Print(printJob.Host, printJob.Message)
	}
}

func Print(host, msg string) {
	redStart := "\033[0;31m"
	redEnd := "\033[0m"
	colourHost := redStart + host + redEnd
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf("%s: %s", colourHost, msg)
	} else {
		fmt.Printf("%s: %s", host, msg)
	}
}

func RunCommand(host, command string) (string, error) {
	sshCmd := []string{
		"-o",
		"StrictHostKeyChecking=no",
		host,
		command,
	}
	cmd := exec.Command("ssh", sshCmd...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
