DESTDIR ?=
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
  PREFIX ?= /usr/local
  INSTALL_TARGET = install-macos
else
  PREFIX ?= /usr
  INSTALL_TARGET = install-linux
endif
BINDIR = $(PREFIX)/bin
PWD = $(shell pwd)

.PHONY: all samples clean install-bin install devinstall uninstall

all: cmd/battlestarc/battlestarc

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
	make -C life

clean:
	(cd cmd/battlestarc; go clean)
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
	make -C life clean

cmd/battlestarc/battlestarc:
	(cd cmd/battlestarc; go build)

install: $(INSTALL_TARGET) cmd/battlestarc/battlestarc

install-macos: cmd/battlestarc/battlestarc
	cp -v "$(PWD)/scripts/bts.sh" "$(DESTDIR)$(BINDIR)/bts"
	cp -v "$(PWD)/scripts/build.sh" "$(DESTDIR)$(BINDIR)/btsbuild"
	cp -v "$(PWD)/cmd/battlestarc/battlestarc" "$(DESTDIR)$(BINDIR)/battlestarc"
	cp -v "$(PWD)/scripts/com2bts.py" "$(DESTDIR)$(BINDIR)/com2bts"
	chmod +x "$(DESTDIR)$(BINDIR)/bts"
	chmod +x "$(DESTDIR)$(BINDIR)/btsbuild"
	chmod +x "$(DESTDIR)$(BINDIR)/battlestarc"
	chmod +x "$(DESTDIR)$(BINDIR)/com2bts"

install-linux: cmd/battlestarc/battlestarc
	install -Dm755 "$(PWD)/scripts/bts.sh" "$(DESTDIR)$(BINDIR)/bts"
	install -Dm755 "$(PWD)/scripts/build.sh" "$(DESTDIR)$(BINDIR)/btsbuild"
	install -Dm755 "$(PWD)/cmd/battlestarc/battlestarc" "$(DESTDIR)$(BINDIR)/battlestarc"
	install -Dm755 "$(PWD)/scripts/com2bts.py" "$(DESTDIR)$(BINDIR)/com2bts"

install-dev: devinstall

devinstall: cmd/battlestarc/battlestarc
	ln -sf $(PWD)/cmd/battlestarc/battlestarc $(BINDIR)/battlestarc
	ln -sf $(PWD)/scripts/bts.sh $(BINDIR)/bts
	ln -sf $(PWD)/scripts/build.sh $(BINDIR)/btsbuild
	ln -sf $(PWD)/scripts/com2bts.py $(BINDIR)/com2bts
	chmod a+rx $(PWD)/cmd/battlestarc/battlestarc
	chmod a+rx $(PWD)/scripts/bts.sh
	chmod a+rx $(PWD)/scripts/build.sh
	chmod a+rx $(PWD)/scripts/com2bts.py

uninstall:
	rm -f "$(DESTDIR)$(BINDIR)/bts"
	rm -f "$(DESTDIR)$(BINDIR)/btsbuild"
	rm -f "$(DESTDIR)$(BINDIR)/battlestarc"
	rm -f "$(DESTDIR)$(BINDIR)/com2bts"
