all: battlestarc

full: clean battlestarc
	make -C test_bits
	make -C samples

full_clean: clean
	make -C test_bits clean
	make -C samples clean

battlestarc:
	@rm -f battlestar # Make sure only "battlestarc" is present, not "battlestar"
	go build -o battlestarc

clean:
	rm -f battlestarc
