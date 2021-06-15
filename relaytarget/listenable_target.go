package relaytarget

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

type listenableRelayTarget struct {
	addr      string
	conn      net.Conn
	client    net.Conn
	tlsConfig *tls.Config
}

func NewListenableRelayTarget(port int, tlsConfig *tls.Config) *listenableRelayTarget {
	return &listenableRelayTarget{
		addr:      fmt.Sprintf(":%d", port),
		tlsConfig: tlsConfig,
	}
}

func (t *listenableRelayTarget) Prepare() error {
	log.Printf("target server is listening on address %s\n", t.addr)
	l, err := net.Listen("tcp4", t.addr)
	if err != nil {
		return err
	}
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		if t.client != nil {
			t.client.Close()
		}
		t.client = conn
	}()
	return nil
}

func (t *listenableRelayTarget) Dial() error {
	return nil
}

func (t *listenableRelayTarget) Close() error {
	err := t.conn.Close()
	if err != nil {
		return err
	}
	err = t.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *listenableRelayTarget) Conn() net.Conn {
	return t.client
}
