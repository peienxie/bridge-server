package relaytarget

import (
	"crypto/tls"
	"net"
)

type relayTarget struct {
	addr      string
	conn      net.Conn
	tlsConfig *tls.Config
}

func NewRelayTarget(addr string, tlsConfig *tls.Config) *relayTarget {
	return &relayTarget{
		addr:      addr,
		tlsConfig: tlsConfig,
	}
}

func (t *relayTarget) Prepare() error {
	return nil
}

func (t *relayTarget) Dial() error {
	var conn net.Conn
	var err error
	if t.tlsConfig != nil {
		conn, err = tls.Dial("tcp4", t.addr, t.tlsConfig)
	} else {
		conn, err = net.Dial("tcp4", t.addr)
	}
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *relayTarget) Close() error {
	err := t.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *relayTarget) Conn() net.Conn {
	return t.conn
}
