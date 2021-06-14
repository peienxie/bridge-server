package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	SERVER_PORT = 8088
)

func main() {
	err := runClient(fmt.Sprintf(":%d", SERVER_PORT))
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
