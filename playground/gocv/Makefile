
.PHONY: ⚙️

all: ⚙️ build-webcam

build-webcam: ⚙️
	go build -o bin/webcam cmd/webcam/main.go
	codesign -f -s - bin/webcam
	bin/webcam -h

reset-camera:
	tccutil reset Camera
