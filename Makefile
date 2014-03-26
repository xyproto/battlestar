DESTDIR ?= 
PREFIX ?= /usr
BINDIR = $(PREFIX)/bin
PWD = $(shell pwd)

all: clean battlestarc

full: clean battlestarc
	make -C helloworld
	make -C test_bits
	make -C test_bits2
	make -C samples
	make -C kernel

full_clean:
	make -C helloworld clean
	make -C test_bits clean
	make -C test_bits2 clean
	make -C samples clean
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
	install -Dm755 scripts/btstool.sh "$(DESTDIR)$(BINDIR)/bts"
	install -Dm755 scripts/btsbuild.sh "$(DESTDIR)$(BINDIR)/btsbuild"
	[ -e $(DESTDIR)$(BINDIR)/objdump ] && install -Dm755 scripts/disasm.sh "$(DESTDIR)$(BINDIR)/disasm"
	install -Dm755 battlestarc "$(DESTDIR)$(BINDIR)/battlestarc"

install: install-bin

devinstall: battlestarc
	ln -sf $(PWD)/battlestarc /usr/bin/battlestarc
	ln -sf $(PWD)/scripts/btstool.sh /usr/bin/bts
	ln -sf $(PWD)/scripts/build.sh /usr/bin/btsbuild

uninstall:
	rm -f "$(DESTDIR)$(BINDIR)/bts"
	rm -f "$(DESTDIR)$(BINDIR)/disasm"
	rm -f "$(DESTDIR)$(BINDIR)/battlestarc"
