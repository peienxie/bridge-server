package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"
	"tcprelay"
	"tcprelay/relaytarget"

	"github.com/alyu/configparser"
)

type Config struct {
	MiddleServerPort    int
	SecuredMiddleServer bool
	TargetServerIp      string
	TargetServerPort    int
	SecuredTargetServer bool
}

func (c Config) String() string {
	return fmt.Sprintf("MiddleServerPort: %d\nSecuredMiddleServer: %v\nTargetServerIp: %s\nTargetServerPort: %d\nSecuredTargetServer: %v\n",
		c.MiddleServerPort,
		c.SecuredMiddleServer,
		c.TargetServerIp,
		c.TargetServerPort,
		c.SecuredTargetServer,
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

func setupMiddleServerTLSConfig() *tls.Config {
	if !appConfig.SecuredMiddleServer {
		return nil
	}

	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatal(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func setupTargetServerTLSConfig() *tls.Config {
	if !appConfig.SecuredTargetServer {
		return nil
	}

	return &tls.Config{InsecureSkipVerify: true}
}

func main() {

	s := tcprelay.NewTcpRelayServer(
		appConfig.MiddleServerPort,
		relaytarget.NewRelayTarget(
			fmt.Sprintf("%s:%d", appConfig.TargetServerIp, appConfig.TargetServerPort),
			setupTargetServerTLSConfig(),
		),
		// relaytarget.NewListenableRelayTarget(8089),
		setupMiddleServerTLSConfig(),
	)
	s.Listen()

	done := make(chan bool, 1)
	<-done
}
