package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sevlyar/go-daemon"
	st "github.com/starryalley/go-servotester/pkg/common"
	"gobot.io/x/gobot/platforms/raspi"
)

//Version of this tool
var Version = "dev"

var rpi *raspi.Adaptor

const programName = "pi-servo-httpd"

func servoHandler(w http.ResponseWriter, r *http.Request) {
	// pin
	keys, ok := r.URL.Query()["pin"]
	if !ok || len(keys[0]) <= 0 {
		http.Error(w, "param 'pin' is missing", http.StatusBadRequest)
		return
	}
	pin := keys[0]

	// position
	keys, ok = r.URL.Query()["pos"]
	if !ok || len(keys[0]) <= 0 {
		http.Error(w, "param 'pos' is missing", http.StatusBadRequest)
		return
	}
	pos, err := strconv.ParseFloat(keys[0], 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Setting pin %v servo position: %.1f", pin, pos)
	p, err := st.CreateServoPacket(pin, pos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// set servo position
	if err = st.UpdatePWM(rpi, *p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func initRaspberryPi(period uint32) (rpi *raspi.Adaptor) {
	rpi = raspi.NewAdaptor()
	rpi.PiBlasterPeriod = period // pi-blaster set to 50Hz
	log.Printf("pi-blaster period set to %v nanoseconds\n", period)
	return
}

func main() {
	log.Printf("%s version:%s\n", programName, Version)
	var port = flag.Int("port", 8080, "Port. Default: 8080")
	var blasterPeriod = flag.Uint64("period", 20000000, "pi-blaster current period setting. Default: 20000000")
	flag.Parse()

	context := &daemon.Context{
		PidFileName: fmt.Sprintf("%s.pid", programName),
		PidFilePerm: 0644,
		LogFileName: fmt.Sprintf("/var/log/%s.log", programName),
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
	rpi = initRaspberryPi(uint32(*blasterPeriod))

	// start http server
	http.HandleFunc("/servo", servoHandler)
	log.Printf("server started at :%d...\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
