package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	DEFAULT_SERVER_IP = "127.0.0.1"
	DEFAULT_SERVER_PORT = 18101
)

func main() {
	ip := flag.String("ip", DEFAULT_SERVER_IP, "server IP address")
    port := flag.Int("port", DEFAULT_SERVER_PORT, "server port number")

	err := runClient(fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		log.Fatal(err)
	}
}

func runClient(addr string) error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nCommand to send: ")
		scanner.Scan()
		text := scanner.Text()
		err := sendMessage(addr, text)
		if err != nil {
			log.Println("can't transmit messages:", err)

		}
	}
}

func writeMessage(conn net.Conn, msg string) error {
	if cw, ok := conn.(interface{ CloseWrite() error }); ok {
		defer cw.CloseWrite()
	} else {
		return fmt.Errorf("connection doesn't implement CloseWrite method")
	}
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func sendMessage(addr, msg string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	// send to socket
	err = writeMessage(conn, msg)
	if err != nil {
		return err
	}

	// listen for reply
	buffer := make([]byte, 1024)
	length, err := conn.Read(buffer)
	if err != nil {
		return err
	}
	fmt.Println(string(buffer[:length]))
	return nil
}
