all: clean battlestar

battlestar:
	go build -o battlestarc

clean:
	rm -f battlestarc
