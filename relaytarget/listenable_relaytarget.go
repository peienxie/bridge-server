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

func NewListenableRelayTarget(addr string, tlsConfig *tls.Config) *listenableRelayTarget {
	return &listenableRelayTarget{
		addr:      addr,
		tlsConfig: tlsConfig,
	}
}

func (t *listenableRelayTarget) Prepare() error {
	l, err := net.Listen("tcp4", t.addr)
	if err != nil {
		return err
	}
	log.Printf("target server is listening on address %s\n", t.addr)
	go func() {
		for {
			conn, err := l.Accept()
			log.Printf("target server is connected from %s\n", conn.RemoteAddr().String())
			if err != nil {
				return
			}
			if t.client != nil {
				t.client.Close()
			}
			t.client = conn
		}
	}()
	return nil
}

func (t *listenableRelayTarget) Dial() error {
	if t.client == nil {
		return fmt.Errorf("target server is not ready")
	}
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
