package main

import (
	"tcprelay"
	"tcprelay/relaytarget"
)

var (
	MIDDLE_SERVER_PORT = 8088
)

func main() {
	s := tcprelay.NewTcpRelayServer(
		MIDDLE_SERVER_PORT,
		// relaytarget.NewRelayTarget("18.235.124.214:80"),
		relaytarget.NewListenableRelayTarget(8089),
	)
	s.Listen()

	done := make(chan bool, 1)
	<-done
}
