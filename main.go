package main

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kurtmc/ec2search/search"
	"golang.org/x/crypto/ssh/terminal"
)

var mutex = &sync.Mutex{}

type PrintJob struct {
	Host    string
	Message string
}

func main() {
	var wg sync.WaitGroup
	instances, err := search.ListInstances(os.Args[1])
	n := len(instances)

	if len(os.Args[1:]) < 2 {
		rand.Seed(time.Now().Unix())
		SshMachine(instances[rand.Intn(n)])
	}

	if err != nil {
		panic(err)
	}

	in := make(chan string)

	wg.Add(1)
	go func() {
		for _, i := range instances {
			in <- i
		}
		close(in)
		wg.Done()
	}()

	out := make(chan PrintJob)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for host := range in {
				if host == "" {
					continue
				}
				commandOutput, err := RunCommand(host, os.Args[2])
				if err != nil {
					fmt.Println(err)
					out <- PrintJob{Host: host, Message: fmt.Sprintf("Failed to execute, err: %v", err)}
					continue
				}
				out <- PrintJob{Host: host, Message: commandOutput}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for printJob := range out {
		Print(printJob.Host, printJob.Message)
	}
}

func Print(host, msg string) {
	colour := hashColour(host)
	start := "\033[38;5;" + colour + "m"
	end := "\033[0m"
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf("%s: %s\n", start+host+end, strings.TrimSpace(msg))
	} else {
		fmt.Printf("%s: %s\n", host, strings.TrimSpace(msg))
	}
}

func SshMachine(host string) {
	out, err := exec.Command("which", "ssh").Output()
	sshPath := strings.TrimSpace(string(out))
	if err != nil {
		panic(err)
	}

	sshCmd := []string{
		"ssh",
		"-o",
		"StrictHostKeyChecking=no",
		host,
	}
	syscall.Exec(sshPath, sshCmd, os.Environ())
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

func hashColour(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", (int(h.Sum32())%256)+1)
}
