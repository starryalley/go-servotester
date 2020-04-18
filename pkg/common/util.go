package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"gobot.io/x/gobot/platforms/raspi"
)

// UpdatePWM writes PWM cycle to specific pin stored in p
func UpdatePWM(rpi *raspi.Adaptor, p ServoPacket) error {
	pin, err := rpi.PWMPin(strconv.Itoa(int(p.PinNo)))
	if err != nil {
		return fmt.Errorf("get Pin %v failed: %v", p.PinNo, err)
	}
	err = pin.SetDutyCycle(p.DutyCycle * 1000) //in nanoseconds
	if err != nil {
		return fmt.Errorf("set duty cycle failed: %v", err)
	}
	log.Printf("Writing dc=%v to pin %v\n", p.DutyCycle, p.PinNo)
	return nil
}

// SendServoPacket sends a ServoPacket
func SendServoPacket(conn net.Conn, p ServoPacket) error {
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	err := binary.Write(conn, binary.LittleEndian, p)
	if err != nil {
		return fmt.Errorf("error sending: %v", err)
	}
	//time.Sleep(1 * time.Millisecond)
	return nil
}

// SendServoPosition sends servo position (0~100%) with presumed endpoint (1.0ms ~ 2.0ms)
func SendServoPosition(conn net.Conn, pin string, pos float64) error {
	p, err := CreateServoPacket(pin, pos)
	if err != nil {
		return err
	}
	return SendServoPacket(conn, *p)
}

// SendServoPWM sends servo position as PWM duty cycle (microseconds)
func SendServoPWM(conn net.Conn, pin string, dc uint32) error {
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

	pinNo, err := strconv.Atoi(pin)
	if err != nil {
		return fmt.Errorf("error converting pin: %v", pin)
	}
	if pinNo > 40 || pinNo < 1 {
		return errors.New("error ping name")
	}
	p := ServoPacket{byte(pinNo), dc}
	err = binary.Write(conn, binary.LittleEndian, p)
	if err != nil {
		return fmt.Errorf("error sending: %v", err)
	}
	//time.Sleep(1 * time.Millisecond)
	return nil
}

// ConnectToServer is a helper function to connect to server
func ConnectToServer(host string) (net.Conn, error) {
	d := net.Dialer{Timeout: 1 * time.Second}
	conn, err := d.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	err = conn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		return nil, err
	}

	err = conn.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
