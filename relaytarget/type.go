package relaytarget

import "net"

type TcpRelayTarget interface {
	Prepare() error
	Dial() error
	Close() error
	Conn() net.Conn
}
