## Installation
Based on https://gocv.io/getting-started/macos
```
brew install opencv
brew install pkgconfig
go get -u -d gocv.io/x/gocv
go install gocv.io/x/gocv
```

## Usage
```
go run cmd/affine/main.go
```

## Camera Access
```
go build -o bin/webcam cmd/webcam/main.go
# go build -o bin/affine cmd/affine/main.go
codesign -f -s - bin/*
tccutil reset Camera
bin/webcam
```
