package common

// ServoPacket is a structure used to transfer pin and PWM duty cycle across network
type ServoPacket struct {
	PinNo     byte   // pin number on raspi
	DutyCycle uint32 // in micro-second
}
