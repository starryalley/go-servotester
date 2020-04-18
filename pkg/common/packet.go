package common

import (
	"errors"
	"fmt"
	"strconv"
)

// ServoPacket is a structure used to transfer pin and PWM duty cycle across network
type ServoPacket struct {
	PinNo     byte   // pin number on raspi
	DutyCycle uint32 // in micro-second
}

// CreateServoPacket create a ServoPacket based on pin and position (0~100)
func CreateServoPacket(pin string, pos float64) (*ServoPacket, error) {
	if pos < 0 || pos > 100 {
		return nil, errors.New("out of range")
	}

	// endpoint [1.0ms, 2.0ms]
	dutyCycle := uint32(float64(0.01*pos+1) * 1000) // in microseconds

	pinNo, err := strconv.Atoi(pin)
	if err != nil {
		return nil, fmt.Errorf("error converting pin %v: %v", pin, err)
	}
	if pinNo > 40 || pinNo < 1 {
		return nil, errors.New("error ping name")
	}
	return &ServoPacket{byte(pinNo), dutyCycle}, nil
}
