COMMIT = $(shell git rev-parse --short HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

all:
	go build -ldflags "-X main.Commit $(BRANCH)-$(COMMIT)" 

deb_package: all
	rm -f galera*watchdog*.deb
	mkdir bin
	cp galera_watchdog bin/
	fpm --prefix=/usr --url https://github.com/crahles/galera_watchdog -s dir -t deb -n galera_watchdog -m'christoph@rahles.de' -v $(BRANCH)-$(COMMIT) bin/galera_watchdog
	rm -rf bin/

install: all
	sudo mv galera_watchdog /usr/bin/

uninstall:
	sudo rm /usr/bin/galera_watchdog

clean:
	rm -f galera_watchdog galera*watchdog*.deb
