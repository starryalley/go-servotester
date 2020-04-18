package main

import (
	"encoding/binary"
	"flag"
	"log"
	"net"
	"strconv"
	"time"

	"gobot.io/x/gobot/platforms/raspi"

	"github.com/sevlyar/go-daemon"
	st "github.com/starryalley/go-servotester/pkg/common"
)

//Version of this tool
var Version = "dev"

func main() {
	log.Println("pi-servotesterd version:", Version)
	var addr = flag.String("addr", "", "Address. Default: \"\"")
	var port = flag.Int("port", 6789, "Port. Default: 6789")
	var blasterPeriod = flag.Uint64("period", 20000000, "pi-blaster current period setting. Default: 20000000")
	flag.Parse()

	context := &daemon.Context{
		PidFileName: "pi-servotesterd.pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/pi-servotesterd.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := context.Reborn()
	if err != nil {
		log.Fatal("Unable to run:", err)
	}
	if d != nil {
		return
	}
	defer context.Release()

	// setup raspi
	rpi := initRaspberryPi(uint32(*blasterPeriod))

	// start server
	src := *addr + ":" + strconv.Itoa(*port)
	listener, err := net.Listen("tcp", src)
	if err != nil {
		log.Fatalf("error listening: %v\n", err)
	}
	defer listener.Close()
	log.Printf("servotesterd listening on %s\n", src)

	// listening for incoming connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		go handleRequest(conn, rpi)
	}
}

func handleRequest(conn net.Conn, rpi *raspi.Adaptor) {
	client := conn.RemoteAddr().String()
	log.Printf("%v connected\n", client)

	count := make(map[byte]uint64)
	ch := make(chan st.ServoPacket, 32)
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
					log.Printf("Total processed requests: %v (%v) [+%.f seconds]\n",
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
		var p st.ServoPacket
		err := binary.Read(conn, binary.LittleEndian, &p)
		if err != nil {
			break
		}
		ch <- p
		count[p.PinNo]++
	}

	close(ch)
	log.Printf("%v disconnected\n", client)
}

func processPacket(ch chan st.ServoPacket, rpi *raspi.Adaptor) {
	for {
		p, ok := <-ch
		if !ok {
			break
		}
		st.UpdatePWM(rpi, p)
	}
}

func initRaspberryPi(period uint32) (rpi *raspi.Adaptor) {
	rpi = raspi.NewAdaptor()
	rpi.PiBlasterPeriod = period // pi-blaster set to 50Hz
	log.Printf("pi-blaster period set to %v nanoseconds\n", period)
	return
}
