package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/tcaine/twamp"
)

func main() {
	interval := flag.Int("interval", 1, "Delay between TWAMP-test requests (seconds)")
	count := flag.Int("count", 5, "Number of requests to send (1..2000000000 packets)")
	rapid := flag.Bool("rapid", false, "Send requests rapidly (default count of 5)")
	size := flag.Int("size", 42, "Size of request packets (0..65468 bytes)")
	tos := flag.Int("tos", 0, "IP type-of-service value (0..255)")
	wait := flag.Int("wait", 1, "Maximum wait time after sending final packet (seconds)")
	port := flag.Int("port", 6666, "UDP port to send request packets")
	mode := flag.String("mode", "ping", "Mode of operation (ping, json)")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("No hostname or IP address was specified.")
		os.Exit(1)
	}

	remoteIP := args[0]
	remoteAddr := net.TCPAddr{IP: net.ParseIP(remoteIP), Port: 862}
	remoteServer := remoteAddr.String()

	c := twamp.NewClient()
	connection, err := c.Connect(remoteServer)
	if err != nil {
		log.Fatal(err)
	}

	var ipVersion = 6
	if remoteAddr.IP.To4() != nil {
		ipVersion = 4
	} else {
		ipVersion = 6
	}

	fmt.Println("create session")
	session, err := connection.CreateSession(
		twamp.TwampSessionConfig{
			ReceiverPort: *port,
			SenderPort:   *port,
			Timeout:      *wait,
			Padding:      *size,
			TOS:          *tos,
			IPVersion:    ipVersion,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("create test")

	test, err := session.CreateTest()
	if err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "json":
		results := test.RunX(*count, nil)
		test.FormatJSON(results)
	case "ping":
		test.Ping(*count, *rapid, *interval)
	}

	session.Stop()
	connection.Close()
}
