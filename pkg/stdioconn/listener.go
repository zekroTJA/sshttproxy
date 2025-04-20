package stdioconn

import (
	"log/slog"
	"net"
)

type Listener struct {
	conn      net.Conn
	closeChan chan struct{}
}

func NewListener(conn net.Conn) *Listener {
	return &Listener{
		conn:      conn,
		closeChan: make(chan struct{}),
	}
}

func NewListenerFromStdInOut() *Listener {
	return NewListener(NewConnFromStdInOut())
}

func (t *Listener) Accept() (net.Conn, error) {
	slog.Debug("accepting connection")

	// The HTTP server will accept incomming connections in a loop, so
	// when the first "connection" is accepted, this will block until
	// the listener has been closed. This is totally viable becuase in
	// this use-case, there is only one "connection" expected.
	//
	// If this is not done, the HTTP server will accept indefinetly
	// in a loop because no call to Accept blocks.
	if t.conn == nil {
		slog.Debug("already accepted, blocking ...")
		<-t.closeChan
	}

	conn := t.conn
	t.conn = nil

	return conn, nil
}

func (t *Listener) Close() error {
	close(t.closeChan)
	return nil
}

func (t *Listener) Addr() net.Addr {
	return t.conn.LocalAddr()
}
