GO ?= go
FILE = daemon/hkhass.go

rpi:
	GOOS=linux GOARCH=arm GOARM=6 $(GO) build $(FILE)