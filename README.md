go-servotester is a simple RC servo tester program on Raspberry Pi using the Go programming language.

It uses the following projects:
- gobot (http://gobot.io/) for managing pin on Raspberry Pi platform

- pi-blaster (https://github.com/sarfata/pi-blaster) for generating PWM in raspberry Pi

Use [my fork](https://github.com/starryalley/pi-blaster) where I modify the period to make the PWM frequency at 50Hz, which is for both analog and digital RC servo.

- libui (https://github.com/andlabs/ui) for GUI version of client program.


# Getting started

I run my raspi headlessly so I make this a simple client-server program. The server runs on raspi and listens on a port receving command to control the RC servo. The client program comes with a GUI or a CLI interface, sending servo position command to the server, and raspi will then set the PWM on specific pin to control the RC servo.

![my setup](images/setup.png)

# Build

The project uses go module with Makefile. Simply use make to build:

`make`

Cross compiling for ARMv6:

`GOARM=6 GOARCH=arm GOOS=linux make`

Cross compiling for ARMv7:

`GOARM=7 GOARCH=arm GOOS=linux make`

Note: Raspberry Pi A, A+, B, B+, Zero use ARM6, and Raspberry Pi 2, 3, 4 use ARM7.

For GUI version of client program, you have to manually build it:

`make build_client_gui`

# Install

`make install` will use systemctl to add the server program on the system (raspberry pi or linux only) and start it.

`make uninstall` to revert it.

The server program's log file will be at `/var/log/pi-servotesterd.log`

You can manually run/install from `bin/` of course.

# Usage

Servo running on raspi, using pi-blaster period 20000000ns, listening on port 6789 (default daemon options)

`./pi-servotesterd -period 20000000 -port 6789`

For GUI control, run this on your desktop

`./servotester_gui`

![Screenshot](images/gui.png)


For cli control, run this on your desktop with correct raspi IP and port

`./servotester -addr 192.168.0.123 -port 6789`

You will see:

```
Connecting to 192.168.0.123:6789
Servo Position>
```

Command format:

`exit`: to exit the program

`12=0`: set pin 12 servo position to 0% (0-100% based on 1ms~2ms endpoints)

`12=23.4`: set pin 12 servo position to 23.4%

`12=swing`: do a full range swing on pin 12

`12=center`: center the servo on pin 12


# Known issue

Latest ui lib will cause GUI program to crash once hitting 'Connect'. I have to debug this when I have time.