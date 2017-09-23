# Version
VERSION = `date +%y.%m`

# If unable to grab the version, default to N/A
ifndef VERSION
    VERSION = "n/a"
endif

#
# Makefile options
#


# State the "phony" targets
.PHONY: all clean build install uninstall


all: clean build

build:
	@echo 'Building go_roguelike...'
	@go build -ldflags '-s -w -X main.Version='${VERSION}

clean:
	@echo 'Cleaning...'
	@go clean

install: build
	@echo installing executable file to /usr/bin/go_roguelike
	@sudo cp trackpadctl /usr/bin/go_roguelike

uninstall: clean
	@echo removing executable file from /usr/bin/go_roguelike
	@sudo rm /usr/bin/go_roguelike
