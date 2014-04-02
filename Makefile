DESTDIR ?= 
PREFIX ?= /usr
BINDIR = $(PREFIX)/bin
PWD = $(shell pwd)

all: clean battlestarc

full: clean battlestarc
	make -C helloworld
	make -C samples
	make -C samples64
	make -C samples32
	make -C samples16
	make -C kernel

full_clean:
	make -C helloworld clean
	make -C samples clean
	make -C samples64 clean
	make -C samples32 clean
	make -C samples16 clean
	make -C kernel clean
	make -C . clean

hello: battlestarc
	make -C helloworld

battlestarc:
	@# Make sure only "battlestarc" is present, not "battlestar"
	@rm -f battlestar
	go build -o battlestarc

clean:
	rm -f battlestarc

install-bin: battlestarc
	install -Dm755 "$(PWD)/scripts/btstool.sh" "$(DESTDIR)$(BINDIR)/bts"
	install -Dm755 "$(PWD)/scripts/build.sh" "$(DESTDIR)$(BINDIR)/btsbuild"
	install -Dm755 "$(PWD)/battlestarc" "$(DESTDIR)$(BINDIR)/battlestarc"

install: install-bin

devinstall: battlestarc
	ln -sf $(PWD)/battlestarc /usr/bin/battlestarc
	ln -sf $(PWD)/scripts/btstool.sh /usr/bin/bts
	ln -sf $(PWD)/scripts/build.sh /usr/bin/btsbuild

uninstall:
	rm -f "$(DESTDIR)$(BINDIR)/bts"
	rm -f "$(DESTDIR)$(BINDIR)/disasm"
	rm -f "$(DESTDIR)$(BINDIR)/battlestarc"
