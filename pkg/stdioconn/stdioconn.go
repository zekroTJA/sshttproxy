package stdioconn

import (
	"errors"
	"io"
	"net"
	"os"
	"time"
)

type dummyAddr string

func (d dummyAddr) Network() string {
	return string(d)
}

func (d dummyAddr) String() string {
	return string(d)
}

type StdioConn struct {
	Rx io.ReadCloser
	Tx io.WriteCloser
}

func NewConnFromStdInOut() *StdioConn {
	return &StdioConn{
		Rx: os.Stdin,
		Tx: os.Stdout,
	}
}

func (s *StdioConn) Read(b []byte) (int, error) {
	return s.Rx.Read(b)
}

func (s *StdioConn) Write(b []byte) (int, error) {
	return s.Tx.Write(b)
}

func (s *StdioConn) Close() error {
	return errors.Join(
		s.Rx.Close(),
		s.Tx.Close())
}

func (s *StdioConn) LocalAddr() net.Addr {
	return dummyAddr("local")
}

func (s *StdioConn) RemoteAddr() net.Addr {
	return dummyAddr("remote")
}

func (s *StdioConn) SetDeadline(t time.Time) error {
	return nil
}

func (s *StdioConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (s *StdioConn) SetWriteDeadline(t time.Time) error {
	return nil
}
