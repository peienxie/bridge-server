package tcprelay

import (
	"fmt"
	"io"
	"log"
	"net"
)

type TcpRelayTargetServer interface {
	Prepare() error
	Dial() error
	Close() error
	Conn() net.Conn
}

type tcpRelayServer struct {
	addr   string
	target TcpRelayTargetServer
}

func NewTcpRelayServer(port int, target TcpRelayTargetServer) *tcpRelayServer {
	return &tcpRelayServer{
		addr:   fmt.Sprintf(":%d", port),
		target: target,
	}
}

func (s *tcpRelayServer) Listen() {
	l, err := net.Listen("tcp4", s.addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("server is listening on address %s\n", s.addr)

	for {
		client, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(client, s.target)
	}
}

func handleConnection(client net.Conn, target TcpRelayTargetServer) {
	defer client.Close()

	err := target.Dial()
	if err != nil {
		log.Printf("can't dial target server: %+v\n", err)
		return
	}

	buffer := make([]byte, 4096)
	log.Printf("start transmission\n")
	n, err := CopyBuffer(target.Conn(), client, buffer)
	if err != nil {
		log.Printf("error when send data by client: %+v\n", err)
		return
	}
	log.Printf("tarnsmit data by client: %d %s\n", n, buffer[:n])
		
	n, err = CopyBuffer(client, target.Conn(), buffer)
	if err != nil {
		log.Printf("error when send data back to client: %+v\n", err)
		return
	}
	log.Printf("receive data from target server: %d %s\n", n, buffer[:n])
}

func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in copyBuffer")
	}

	for {
		nr, er := src.Read(buf)
		log.Printf("readed %d, %s\n", nr, buf[:nr])
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
