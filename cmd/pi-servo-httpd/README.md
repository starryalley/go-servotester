# pi-servo-httpd

This is a simple HTTP server exposing 1 endpoint `servo` on a Raspberry Pi.

This program uses [pi-blaster](https://github.com/sarfata/pi-blaster). See [my fork](https://github.com/starryalley/pi-blaster) which modifies cycle time to 50Hz for RC servo.

# Endpoint

`servo` takes 2 query parameters:

- `pin`: the pin number on Raspberry Pi
- `pos`: the servo position. Value is a float from `0.0` to `100.0`, meaning 0% to 100% position of a RC servo.

## Examples:

`http://127.0.0.1:8080/servo?pin=11&pos=15`


# Use

See my other project [Bike Fan Speed](https://github.com/starryalley/BikeFanSpeed) which uses this program.
