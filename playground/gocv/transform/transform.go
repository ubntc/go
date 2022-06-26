package transform

import (
	"image"

	"gocv.io/x/gocv"
)

func ScaleUp(img gocv.Mat, dst *gocv.Mat, scale float64) {
	gocv.Resize(img, dst, image.Point{0, 0}, scale, scale, gocv.InterpolationLinear)
}

func AffineTransform(img gocv.Mat, dst *gocv.Mat, w, h int) {
	rot := gocv.GetRotationMatrix2D(image.Point{w / 2, h / 2}, 90.0, 0.75)
	defer rot.Close()
	gocv.WarpAffine(img, dst, rot, image.Point{w, h})
}
