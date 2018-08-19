package main

import (
	"fmt"
	"github.com/andlabs/ui"
	"net"
	"strconv"
)

func main() {
	// tcp connection to server, and the state
	var conn net.Conn
	connected := false

	// servo endpoints and position, in microseconds
	var lowerEndpoint uint32 = 1000
	var upperEndpoint uint32 = 2000
	const centrePosition = 1520
	var currentPosition uint32 = centrePosition

	ui.Main(func() {
		// servo position
		hboxPos := ui.NewHorizontalBox()
		hboxPos.Append(ui.NewLabel("Servo Position: "), false)
		servoPositionLabel := ui.NewLabel("1520")
		hboxPos.Append(servoPositionLabel, false)
		hboxPos.Append(ui.NewLabel("us"), false)

		// the pin name input and the position slider
		pinInput := ui.NewEntry()
		pinInput.SetText("12")
		slider := ui.NewSlider(0, 100000) // 100000-step resolution

		// the function to update servo position by setting dutyCycle
		updateServoPosition := func(dutyCycle uint32) {
			currentPosition = dutyCycle
			// update label
			servoPositionLabel.SetText(strconv.Itoa(int(dutyCycle)))

			// update slider position
			pos := int(float32(dutyCycle-lowerEndpoint) / float32(upperEndpoint-lowerEndpoint) *
				float32(100000))
			slider.SetValue(pos)

			go SendServoPWM(conn, pinInput.Text(), dutyCycle)
		}

		// slider on change
		slider.OnChanged(func(*ui.Slider) {
			if pinInput.Text() != "" && connected {
				// calculate dutyCycle based on current endpoint
				currentPosition = uint32(float32(slider.Value())/float32(100000)*
					float32(upperEndpoint-lowerEndpoint) + float32(lowerEndpoint))

				// update label
				servoPositionLabel.SetText(strconv.Itoa(int(currentPosition)))
				// send to server
				go SendServoPWM(conn, pinInput.Text(), currentPosition)
			}
		})

		// server ip/port and connect buttons
		statusLabel := ui.NewLabel("Not connected")
		hbox := ui.NewHorizontalBox()
		serverIpInput := ui.NewEntry()
		serverIpInput.SetText("192.168.0.123")
		serverPortInput := ui.NewEntry()
		serverPortInput.SetText("6789")
		connectButton := ui.NewButton("Connect")
		connectButton.OnClicked(func(b *ui.Button) {
			var err error
			if !connected {
				conn, err = ConnectToServer(serverIpInput.Text() + ":" + serverPortInput.Text())
				if err != nil {
					fmt.Println(err)
					statusLabel.SetText(err.Error())
				} else {
					statusLabel.SetText("Connected")
					b.SetText("Disconnect")
					connected = true
					// set initial servo position
					updateServoPosition(currentPosition)
				}
			} else {
				conn.Close()
				statusLabel.SetText("Disconnected")
				b.SetText("Connect")
				connected = false
				// reset currentPosition
				currentPosition = centrePosition
			}
		})
		hbox.Append(serverIpInput, true)
		hbox.Append(serverPortInput, false)
		hbox.Append(connectButton, false)

		// servo endpoints spinboxes
		hboxEP := ui.NewHorizontalBox()
		lowerSpinbox := ui.NewSpinbox(900, 1100)  // conservative: go only 100us beyond
		upperSpinbox := ui.NewSpinbox(1900, 2100) // conservative: go only 100us beyond
		lowerSpinbox.SetValue(int(lowerEndpoint))
		upperSpinbox.SetValue(int(upperEndpoint))
		lowerSpinbox.OnChanged(func(s *ui.Spinbox) {
			if connected {
				lowerEndpoint = uint32(s.Value())
				updateServoPosition(lowerEndpoint)
			}
		})
		upperSpinbox.OnChanged(func(s *ui.Spinbox) {
			if connected {
				upperEndpoint = uint32(s.Value())
				updateServoPosition(upperEndpoint)
			}
		})
		hboxEP.Append(ui.NewLabel("LowerEP:"), false)
		hboxEP.Append(lowerSpinbox, true)
		hboxEP.Append(ui.NewLabel("UpperEP:"), false)
		hboxEP.Append(upperSpinbox, true)

		// servo quick position buttons
		hboxQuickButton := ui.NewHorizontalBox()
		swingButton := ui.NewButton("Swing")
		swingButton.OnClicked(func(b *ui.Button) {
			if connected {
				go func() {
					const step = 20
					for i := lowerEndpoint; i <= upperEndpoint; i += step {
						SendServoPWM(conn, pinInput.Text(), i)
					}
					for i := upperEndpoint; i >= lowerEndpoint; i -= step {
						SendServoPWM(conn, pinInput.Text(), i)
					}
					updateServoPosition(lowerEndpoint)
				}()
			}
		})
		hboxQuickButton.Append(swingButton, false)
		centreButton := ui.NewButton("Centre")
		centreButton.OnClicked(func(b *ui.Button) {
			if connected {
				currentPosition = centrePosition
				updateServoPosition(currentPosition)
			}
		})
		hboxQuickButton.Append(centreButton, false)

		// main vertical layout
		vbox := ui.NewVerticalBox()
		vbox.Append(statusLabel, false)
		vbox.Append(hbox, false)
		vbox.Append(ui.NewHorizontalSeparator(), false)
		vbox.Append(ui.NewLabel("Pin Name on Raspi"), false)
		vbox.Append(pinInput, false)
		vbox.Append(hboxPos, false)
		vbox.Append(slider, false)
		vbox.Append(hboxEP, false)
		vbox.Append(hboxQuickButton, false)
		vbox.SetPadded(true)

		// main window
		window := ui.NewWindow("ServoTester Controller", 400, 250, true)
		window.SetMargined(true)
		window.SetChild(vbox)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			if connected {
				conn.Close()
			}
			return true
		})
		window.Show()
	})
}
