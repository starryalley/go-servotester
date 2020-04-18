ifeq ($(origin VERSION), undefined)
	VERSION := $(shell git rev-parse --short HEAD)
endif

PACKAGENAME = main

.PHONY: all

CMD := pi-servotesterd servotester

all: build

build:
	@echo "Version:$(VERSION)"
	for target in $(CMD); do \
		$(BUILD_ENV_FLAGS) go build -v -o bin/$$target -ldflags "-X $(PACKAGENAME).Version=$(VERSION)" ./cmd/$$target; \
	done

test:
	go test ./...

build_client_gui:
	go build -v -o bin/servotester_gui "-X $(PACKAGENAME).Version=$(VERSION)" ./cmd/servotester_gui

install:
	cp bin/pi-servotesterd /usr/sbin
	cp pi-servotester.service /lib/systemd/system
	systemctl enable pi-servotester
	systemctl start pi-servotester

uninstall:
	systemctl stop pi-servotester
	systemctl disable pi-servotester
	rm -f /lib/systemd/system/pi-servotester.service
	rm -f /usr/sbin/pi-servotesterd

clean:
	rm -rf ./bin/*

