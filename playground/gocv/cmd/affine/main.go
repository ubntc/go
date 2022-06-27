package main

import (
	"flag"
	"log"

	"github.com/pkg/errors"
	"github.com/ubntc/go/playground/gocv/transform"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/cuda"
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

	w := int(webcam.Get(gocv.VideoCaptureFrameWidth))
	h := int(webcam.Get(gocv.VideoCaptureFrameHeight))

	wnd := gocv.NewWindow("Webcam")
	defer wnd.Close()

	img := gocv.NewMat()
	uimg := cuda.NewGpuMatFromMat(img)
	scaled := gocv.NewMat()
	uscaled := cuda.NewGpuMat()
	udst := cuda.NewGpuMat()
	dst := gocv.NewMat()

	defer img.Close()
	defer uimg.Close()
	defer scaled.Close()
	defer uscaled.Close()
	defer dst.Close()
	defer udst.Close()

	log.Printf("using webcam: %v (%v x %v)\n", dev, w, h)
	for {
		if ok := webcam.Read(&img); !ok {
			return
		}
		if img.Empty() {
			continue
		}

		scale := 5
		ws, hs := w*scale, h*scale
		transform.GpuScaleUp(uimg, &uscaled, float64(scale))

		transform.GpuAffineTransform(uscaled, &udst, ws, hs)
		wnd.IMShow(dst)

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
