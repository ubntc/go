package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"
)

const (
	KEYCODE_NONE = -1
	KEYCODE_ESC  = 27
	KEYCODE_Q    = 113
)

var deviceID = flag.String("deviceID", "0", "device ID passewd to OpenCV for accessing the webcam")

func main() {
	flag.Parse()

	dev := *deviceID

	webcam, err := gocv.OpenVideoCapture(dev)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "OpenVideoCapture: %v", dev))
	}
	defer webcam.Close()

	wnd := gocv.NewWindow("Webcam")
	defer wnd.Close()

	img := gocv.NewMat()
	defer img.Close()

	fmt.Printf("using webcam: %v\n", dev)
	for {
		if ok := webcam.Read(&img); !ok {
			return
		}
		if img.Empty() {
			continue
		}

		wnd.IMShow(img)
		switch key := wnd.WaitKey(1); key {
		case KEYCODE_NONE:
			continue
		case KEYCODE_ESC, KEYCODE_Q:
			return
		default:
			log.Println("unsupported keycode", key)
		}
	}
}
