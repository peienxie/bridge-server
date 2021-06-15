package main

import (
	"fmt"
	"log"
	"strconv"
	"tcprelay"
	"tcprelay/relaytarget"

	"github.com/alyu/configparser"
)

type Config struct {
	MiddleServerPort int
	TargetServerIp   string
	TargetServerPort int
}

func (c Config) String() string {
	return fmt.Sprintf("MiddleServerPort: %d\nTargetServerIp: %s\nTargetServerPort: %d\n",
		c.MiddleServerPort,
		c.TargetServerIp,
		c.TargetServerPort,
	)
}

var appConfig Config

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	config, err := configparser.Read("config.ini")
	checkErr(err)
	section, err := config.Section("CONFIG")
	checkErr(err)

	appConfig.MiddleServerPort, err = strconv.Atoi(section.ValueOf("MiddleServerPort"))
	checkErr(err)

	appConfig.TargetServerIp = section.ValueOf("TargetServerIp")

	appConfig.TargetServerPort, err = strconv.Atoi(section.ValueOf("TargetServerPort"))
	checkErr(err)

	fmt.Printf("using configuration below:\n%s\n", appConfig)
}

func main() {
	s := tcprelay.NewTcpRelayServer(
		appConfig.MiddleServerPort,
		relaytarget.NewRelayTarget(fmt.Sprintf("%s:%d", appConfig.TargetServerIp, appConfig.TargetServerPort)),
		// relaytarget.NewListenableRelayTarget(8089),
	)
	s.Listen()

	done := make(chan bool, 1)
	<-done
}
