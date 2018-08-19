.PHONY: build build_arm6 build_arm7 build_client clean

all: build_arm6 build_client 

build_arm6:
	GOARM=6 GOARCH=arm GOOS=linux go build -o pi-servotesterd server.go packet.go

build_arm7:
	GOARM=7 GOARCH=arm GOOS=linux go build -o pi-servotesterd server.go packet.go

build_client:
	go build -o servotester_cli util.go packet.go client.go
	go build -o servotester util.go packet.go client_gui.go

clean:
	rm servotester servotester_cli pi-servotesterd

