package tcprelay

import (
	"bufio"
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
	log.Printf("middle server is listening on address %s\n", s.addr)
	err = s.target.Prepare()
	if err != nil {
		log.Printf("target server is not ready: %+v", err)
		return
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("client connected from %s\n", client.RemoteAddr().String())
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

	err = copy(target.Conn(), client)
	if err != nil {
		log.Printf("error when send data by client: %+v\n", err)
		return
	}

	err = copy(client, target.Conn())
	if err != nil {
		log.Printf("error when send data back to client: %+v\n", err)
		return
	}
}

func copy(dst net.Conn, src net.Conn) (err error) {
	r := bufio.NewReader(src)
	w := bufio.NewWriter(dst)
	buf := make([]byte, 4096)

	buf[0], err = r.ReadByte()
	if err != nil {
		return err
	}
	err = w.WriteByte(buf[0])
	if err != nil {
		return err
	}

	for r.Buffered() > 0 {
		nr, er := r.Read(buf[:])
		if nr > 0 {
			nw, ew := w.Write(buf[:nr])
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
	w.Flush()
	return err
}
