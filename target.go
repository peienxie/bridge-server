package tcprelay

import (
	"net"
)


type relayTarget struct {
	addr string
	conn net.Conn
}


func NewRelayTarget(addr string) *relayTarget {
	return &relayTarget{
		addr : addr,
	}
}

func (t *relayTarget) Prepare() error {
	return nil
}

func (t *relayTarget) Dial() error {
	conn, err := net.Dial("tcp4", t.addr)
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