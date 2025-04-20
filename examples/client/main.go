package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"
	"syscall"

	"github.com/zekrotja/sshttproxy/pkg/client"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Printf("usage: %s [username@]address[:port]\n", os.Args[0])
		os.Exit(1)
	}

	username, host, port, err := splitAddress(os.Args[1])
	if err != nil {
		fmt.Printf("failed getting username: %s\n", err)
		os.Exit(1)
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("failed getting password: %s\n", err)
		os.Exit(1)
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(bytePassword)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", host+":"+port, sshConfig)
	if err != nil {
		fmt.Printf("failed connecting to SSH server: %s\n", err)
		os.Exit(1)
	}
	defer sshClient.Close()

	c, err := client.New(sshClient, "http")
	if err != nil {
		fmt.Printf("failed creating proxy client: %s\n", err)
		os.Exit(1)
	}
	defer c.Close()

	r, _ := http.NewRequest("GET", "http://server/foo?bar=baz", nil)
	r.Header.Set("X-Test", "foobar")
	resp, err := c.Do(r)
	if err != nil {
		fmt.Printf("request failed: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body:\n%s\n", body)
}

func splitAddress(address string) (username, host, port string, err error) {
	split := strings.SplitN(address, "@", 2)
	if len(split) == 1 {
		u, err := user.Current()
		if err != nil {
			return "", "", "", err
		}
		username = u.Username
	} else {
		username = split[0]
		address = split[1]
	}

	split = strings.SplitN(address, ":", 2)
	if len(split) == 1 {
		host = split[0]
		port = "22"
	} else {
		host = split[0]
		port = split[1]
	}

	return username, host, port, nil
}
