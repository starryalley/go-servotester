package main

import (
	"flag"
	"log"
	"strconv"

	st "github.com/starryalley/go-servotester/pkg/common"
)

//Version of this tool
var Version = "dev"

func main() {
	var addr = flag.String("addr", "127.0.0.1", "Address. Default: 127.0.0.1")
	var port = flag.Int("port", 6789, "Port. Default: 6789")
	var pin = flag.String("pin", "11", "Pin number. Default: 11")
	var position = flag.Int("position", 50, "Servo position. 0-100. Default: 50")
	flag.Parse()

	host := *addr + ":" + strconv.Itoa(*port)

	conn, err := st.ConnectToServer(host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("Connecting to %v, setting pin %s to position %v. Tool version: %s\n", host, *pin, *position, Version)
	err = st.SendServoPosition(conn, *pin, float64(*position))
	if err != nil {
		log.Fatal(err)
	}
}
