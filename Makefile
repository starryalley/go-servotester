ifeq ($(origin VERSION), undefined)
	VERSION := $(shell git rev-parse --short HEAD)
endif

PACKAGENAME = main

.PHONY: all

CMD := pi-servotesterd servotester servo_control pi-servo-httpd

all: build

build:
	@echo "Version:$(VERSION)"
	for target in $(CMD); do \
		$(BUILD_ENV_FLAGS) go build -v -o bin/$$target -ldflags "-X $(PACKAGENAME).Version=$(VERSION)" ./cmd/$$target; \
	done

test:
	go test ./...

build_client_gui:
	go build -v -o bin/servotester_gui ./cmd/servotester_gui

install:
	cp bin/pi-servotesterd bin/pi-servo-httpd /usr/sbin
	cp pi-servotester.service pi-servo-http.service /lib/systemd/system
	systemctl enable pi-servotester pi-servo-http
	systemctl start pi-servotester pi-servo-http

uninstall:
	systemctl stop pi-servotester pi-servo-http
	systemctl disable pi-servotester pi-servo-http
	rm -f /lib/systemd/system/pi-servotester.service /lib/systemd/system/pi-servo-http.service
	rm -f /usr/sbin/pi-servotesterd /usr/sbin/pi-servo-httpd

clean:
	rm -rf ./bin/*

