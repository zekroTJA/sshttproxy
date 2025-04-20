package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/zekrotja/sshttproxy/pkg/stdioconn"
	"golang.org/x/crypto/ssh"
)

// Client wraps an SSH client and session to perform HTTP requests
// trhough the SSH subsystem proxy.
type Client struct {
	sshClient *ssh.Client
	session   *ssh.Session
	conn      net.Conn
}

// New creates a new Client with the given SSH client and subsystem name.
func New(sshClient *ssh.Client, subsystemName string) (t *Client, err error) {
	if sshClient == nil {
		return nil, fmt.Errorf("sshClient can not be null")
	}

	t = &Client{}

	t.session, err = sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	stdIn, err := t.session.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdOut, err := t.session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = t.session.RequestSubsystem(subsystemName)
	if err != nil {
		return nil, err
	}

	t.conn = &stdioconn.StdioConn{Rx: io.NopCloser(stdOut), Tx: stdIn}

	return t, nil
}

// Close closes the ssh session and connection.
func (t *Client) Close() error {
	return errors.Join(
		t.conn.Close(),
		t.session.Close(),
	)
}

// Do performs the given http request through the proxy subsystem on the host server.
//
// When the request was successful, the response object is reutrned.
func (t *Client) Do(r *http.Request) (*http.Response, error) {
	err := r.Write(t.conn)
	if err != nil {
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(t.conn), r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
