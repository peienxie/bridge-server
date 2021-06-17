package tcprelay

import (
	"bufio"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"time"
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
	tlsCfg *tls.Config
}

func NewTcpRelayServer(port int, target TcpRelayTargetServer, tlsCfg *tls.Config) *tcpRelayServer {
	return &tcpRelayServer{
		addr:   fmt.Sprintf(":%d", port),
		target: target,
		tlsCfg: tlsCfg,
	}
}

func (s *tcpRelayServer) listener() (l net.Listener, err error) {
	if s.tlsCfg != nil {
		l, err = tls.Listen("tcp", s.addr, s.tlsCfg)
		log.Printf("using TLS connection\n")
	} else {
		l, err = net.Listen("tcp4", s.addr)
	}
	if err != nil {
		return nil, err
	}
	log.Printf("middle server is listening on address %s\n", s.addr)
	return l, nil
}

func (s *tcpRelayServer) Listen() {
	l, err := s.listener()
	if err != nil {
		log.Fatal(err)
	}
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
	defer func() {
		log.Println("closing client connection")
		client.Close()
	}()

	err := target.Dial()
	if err != nil {
		log.Printf("can't dial target server: %+v\n", err)
		return
	}
	log.Println("successfully dial target server", target.Conn().RemoteAddr().String())
	log.Printf("\nstart exchange message\n")

	for {
		log.Printf("\n%s ==========> %s\n", client.RemoteAddr().String(), target.Conn().RemoteAddr().String())
		err = copy(target.Conn(), client)
		if err != nil {
			log.Printf("error when send data by client: %+v\n", err)
			return
		}

		log.Printf("%s <========== %s\n", client.RemoteAddr().String(), target.Conn().RemoteAddr().String())
		err = copy(client, target.Conn())
		if err != nil {
			log.Printf("error when send data back to client: %+v\n", err)
			return
		}
	}
}

func copy(dst net.Conn, src net.Conn) (err error) {
	src.SetReadDeadline(time.Now().Add(time.Second * 10))

	r := bufio.NewReader(src)
	w := bufio.NewWriter(dst)
	buf := make([]byte, 1024)
	data := make([]byte, 0)

	buf[0], err = r.ReadByte()
	if err != nil {
		return fmt.Errorf("read first byte error: %+v\n", err)
	}
	data = append(data, buf[0])
	err = w.WriteByte(buf[0])
	if err != nil {
		return fmt.Errorf("write first byte error: %+v\n", err)
	}

	for r.Buffered() > 0 {
		nr, er := r.Read(buf[:])
		data = append(data, buf[:nr]...)
		if nr > 0 {
			nw, ew := w.Write(buf[:nr])
			if ew != nil {
				err = fmt.Errorf("write data error: %+v\n", ew)
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = fmt.Errorf("read data error: %+v\n", er)
			}
			break
		}
	}
	log.Printf("transmitted packet length:%d\n%s\n", len(data), hex.EncodeToString(data))
	w.Flush()
	return err
}
