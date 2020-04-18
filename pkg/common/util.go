package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

// SendServoPosition sends servo position (0~100%) with presumed endpoint (1.0ms ~ 2.0ms)
func SendServoPosition(conn net.Conn, pin string, pos float64) (err error) {
	if pos < 0 || pos > 100 {
		return errors.New("out of range")
	}
	var p ServoPacket
	// endpoint [1.0ms, 2.0ms]
	p.DutyCycle = uint32(float64(0.01*pos+1) * 1000) // in microseconds

	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

	pinNo, err := strconv.Atoi(pin)
	if err != nil {
		fmt.Printf("Error converting pin: %v", pin)
		return
	}
	if pinNo > 40 || pinNo < 1 {
		errors.New("error ping name")
		return
	}
	p.PinNo = byte(pinNo)
	err = binary.Write(conn, binary.LittleEndian, p)
	if err != nil {
		fmt.Printf("Error sending: %v\n", err)
		return
	}
	//time.Sleep(1 * time.Millisecond)
	return
}

// SendServoPWM sends servo position as PWM duty cycle (microseconds)
func SendServoPWM(conn net.Conn, pin string, dc uint32) (err error) {
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

	pinNo, err := strconv.Atoi(pin)
	if err != nil {
		fmt.Printf("Error converting pin: %v", pin)
		return
	}
	if pinNo > 40 || pinNo < 1 {
		errors.New("error ping name")
		return
	}
	p := ServoPacket{byte(pinNo), dc}
	err = binary.Write(conn, binary.LittleEndian, p)
	if err != nil {
		fmt.Printf("Error sending: %v\n", err)
		return
	}
	//time.Sleep(1 * time.Millisecond)
	return
}

// ConnectToServer is a helper function to connect to server
func ConnectToServer(host string) (conn net.Conn, err error) {
	d := net.Dialer{Timeout: 1 * time.Second}
	conn, err = d.Dial("tcp", host)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = conn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = conn.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
