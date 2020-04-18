package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
    "log"
	"strconv"
	"strings"

    . "github.com/starryalley/go-servotester/pkg/common"
)

//Version of this tool
var Version = "dev"

func main() {
    log.Println("servotester version:", Version)
	var addr = flag.String("addr", "127.0.0.1", "Address. Default: 127.0.0.1")
	var port = flag.Int("port", 6789, "Port. Default: 6789")
	flag.Parse()

	host := *addr + ":" + strconv.Itoa(*port)
	fmt.Printf("Connecting to %v\n", host)

	conn, err := ConnectToServer(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		fmt.Print("Servo Position> ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
			continue
		}
		input = strings.TrimRight(input, "\n")
		// exit command
		if input == "exit" {
			break
		}
		// split by =
		tokens := strings.Split(input, "=")
		if len(tokens) != 2 {
			fmt.Println("Error input format. Use pinname=cmd")
			continue
		}
		fmt.Printf("Pin:%v Cmd:%v\n", tokens[0], tokens[1])

		// special command: swing
		if tokens[1] == "swing" {
			go func() {
				for i := 0.0; i < 100.0; i += 2 {
					SendServoPosition(conn, tokens[0], float64(i))
				}
				for i := 100.0; i >= 0.0; i -= 2 {
					SendServoPosition(conn, tokens[0], float64(i))
				}
			}()
			continue
		}
		if tokens[1] == "center" {
			go SendServoPosition(conn, tokens[0], 50)
			continue
		}

		// normal floating number input
		f, err := strconv.ParseFloat(tokens[1], 64)
		if err != nil {
			fmt.Printf("Error parsing input: %v\n", err)
			continue
		}
		go SendServoPosition(conn, tokens[0], f)
	}
}
