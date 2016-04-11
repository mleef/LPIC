# Simple building of lpic binary
GO=go

all: master

master:	
	$(GO) build -o lpic master.go

clean:
	rm -f lpic