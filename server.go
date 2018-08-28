package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"net"
	"strconv"
	"time"
)

func main() {
	var addr = flag.String("addr", "", "Address. Default: \"\"")
	var port = flag.Int("port", 6789, "Port. Default: 6789")
	var blasterPeriod = flag.Uint64("period", 10000000, "pi-blaster current period setting. Default: 10000000")
	flag.Parse()

	// setup raspi
	rpi := initRaspberryPi(uint32(*blasterPeriod))

	// start server
	src := *addr + ":" + strconv.Itoa(*port)
	listener, err := net.Listen("tcp", src)
	if err != nil {
		log.Fatalf("error listening: %v\n", err)
	}
	defer listener.Close()
	fmt.Printf("servotesterd listening on %s\n", src)

	// listening for incoming connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		go handleRequest(conn, rpi)
	}
}

func handleRequest(conn net.Conn, rpi *raspi.Adaptor) {
	client := conn.RemoteAddr().String()
	fmt.Printf("%v connected\n", client)

	count := make(map[byte]uint64)
	ch := make(chan ServoPacket, 32)
	go processPacket(ch, rpi)

	// reporter timer which will report processed request count every 5 seconds (if there is any change)
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	defer close(quit)
	go func() {
		var lastCount uint64
		lastTime := time.Now()
		// waits for timer tick
		for {
			select {
			case <-ticker.C:
				// calculate total request count
				var currentCount uint64
				for _, v := range count {
					currentCount += v
				}
				// report received count
				if lastCount != currentCount {
					fmt.Printf("Total processed requests: %v (%v) [+%.f seconds]\n",
						currentCount, count, time.Since(lastTime).Seconds())
					lastCount, lastTime = currentCount, time.Now()
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// main TCP stream receiver loop
	for {
		var p ServoPacket
		err := binary.Read(conn, binary.LittleEndian, &p)
		if err != nil {
			break
		}
		ch <- p
		count[p.PinNo]++
	}

	close(ch)
	fmt.Printf("%v disconnected\n", client)
}

func processPacket(ch chan ServoPacket, rpi *raspi.Adaptor) {
	for {
		p, ok := <-ch
		if !ok {
			break
		}
		updatePWM(rpi, p)
	}
}

func initRaspberryPi(period uint32) (rpi *raspi.Adaptor) {
	rpi = raspi.NewAdaptor()
	rpi.PiBlasterPeriod = period // pi-blaster set to 50Hz
	fmt.Printf("pi-blaster period set to %v nanoseconds\n", period)
	return
}

// duty cycle is in nanoseconds
func updatePWM(rpi *raspi.Adaptor, p ServoPacket) (err error) {
	pin, err := rpi.PWMPin(strconv.Itoa(int(p.PinNo)))
	if err != nil {
		fmt.Printf("Get Pin %v failed: %v\n", p.PinNo, err)
		return
	}
	err = pin.SetDutyCycle(p.DutyCycle * 1000) //in nanoseconds
	if err != nil {
		fmt.Println("Set duty cycle failed:", err)
		return
	}
	// fmt.Printf("Writing dc=%v to pin %v\n", p.DutyCycle, p.PinNo)
	return nil
}
