package main

// simple struct used to transfer across network
type ServoPacket struct {
	PinNo     byte   // pin number on raspi
	DutyCycle uint32 // in micro-second
}
