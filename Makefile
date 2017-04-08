DESTDIR ?= 
PREFIX ?= /usr
BINDIR = $(PREFIX)/bin
PWD = $(shell pwd)

.PHONY: all samples clean install-bin install devinstall uninstall

all: src/battlestarc

samples:
	make -C helloworld
	make -C samples
	make -C samples64
	make -C samples32
	make -C samples16
	make -C kernel/simple
	make -C kernel/with_c
	make -C kernel/reverse_string
	make -C bottles99
	make -C fibonacci
	make -C sdl2
	make -C life

clean:
	make -C src clean
	make -C helloworld clean
	make -C samples clean
	make -C samples64 clean
	make -C samples32 clean
	make -C samples16 clean
	make -C kernel/simple clean
	make -C kernel/with_c clean
	make -C kernel/reverse_string clean
	make -C bottles99 clean
	make -C fibonacci clean
	make -C sdl2 clean
	make -C life clean

src/battlestarc:
	make -C src

install-linux: src/battlestarc
	install -Dm755 "$(PWD)/scripts/bts.sh" "$(DESTDIR)$(BINDIR)/bts"
	install -Dm755 "$(PWD)/scripts/build.sh" "$(DESTDIR)$(BINDIR)/btsbuild"
	install -Dm755 "$(PWD)/src/battlestarc" "$(DESTDIR)$(BINDIR)/battlestarc"
	install -Dm755 "$(PWD)/scripts/com2bts.py" "$(DESTDIR)$(BINDIR)/com2bts"

install-osx: uninstall src/battlestarc
	cp -v "$(PWD)/scripts/bts.sh" "$(DESTDIR)$(BINDIR)/bts"
	cp -v "$(PWD)/scripts/build.sh" "$(DESTDIR)$(BINDIR)/btsbuild"
	cp -v "$(PWD)/src/battlestarc" "$(DESTDIR)$(BINDIR)/battlestarc"
	cp -v "$(PWD)/scripts/com2bts.py" "$(DESTDIR)$(BINDIR)/com2bts"
	chmod +x "$(DESTDIR)$(BINDIR)/bts"
	chmod +x "$(DESTDIR)$(BINDIR)/btsbuild"
	chmod +x "$(DESTDIR)$(BINDIR)/battlestarc"

install:
	@# TODO: Add OS X detection
	@echo 'Use "make install-linux" for Linux and "make install-osx" for OS X'

devinstall: src/battlestarc
	ln -sf $(PWD)/src/battlestarc /usr/bin/battlestarc
	ln -sf $(PWD)/scripts/bts.sh /usr/bin/bts
	ln -sf $(PWD)/scripts/build.sh /usr/bin/btsbuild
	ln -sf $(PWD)/scripts/com2bts.py /usr/bin/com2bts
	chmod a+rx $(PWD)/src/battlestarc
	chmod a+rx $(PWD)/scripts/bts.sh
	chmod a+rx $(PWD)/scripts/build.sh
	chmod a+rx $(PWD)/scripts/com2bts.py

uninstall:
	rm -f "$(DESTDIR)$(BINDIR)/bts"
	rm -f "$(DESTDIR)$(BINDIR)/btsbuild"
	rm -f "$(DESTDIR)$(BINDIR)/battlestarc"
	rm -f "$(DESTDIR)$(BINDIR)/com2bts"
