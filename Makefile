DESTDIR ?= 
PREFIX ?= /usr
BINDIR = $(PREFIX)/bin

all: clean battlestarc

full: clean battlestarc
	make -C helloworld
	make -C test_bits
	make -C samples
	make -C kernel

full_clean: clean
	make -C helloworld clean
	make -C test_bits clean
	make -C samples clean
	make -C kernel clean
	make -C osx clean

osx:
	make -C osx

battlestarc:
	@rm -f battlestar # Make sure only "battlestarc" is present, not "battlestar"
	go build -o battlestarc

clean:
	rm -f battlestarc

install-bin: battlestarc
	install -Dm755 scripts/btstool.sh "$(DESTDIR)$(BINDIR)/bts"
	[ -e $(DESTDIR)$(BINDIR)/objdump ] && install -Dm755 scripts/disasm.sh "$(DESTDIR)$(BINDIR)/disasm"
	install -Dm755 battlestarc "$(DESTDIR)$(BINDIR)/battlestarc"

install: install-bin

uninstall:
	rm -f "$(DESTDIR)$(BINDIR)/bts"
	rm -f "$(DESTDIR)$(BINDIR)/disasm"
	rm -f "$(DESTDIR)$(BINDIR)/battlestarc"
