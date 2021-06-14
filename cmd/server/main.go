package main

import (
	"tcprelay"
)

var (
	MIDDLE_SERVER_PORT = 8088
)

func main() {
	s := tcprelay.NewTcpRelayServer(
		MIDDLE_SERVER_PORT,
		tcprelay.NewRelayTarget("18.235.124.214:80"),
	)
	s.Listen()

	done := make(chan bool, 1)
	<-done
}
