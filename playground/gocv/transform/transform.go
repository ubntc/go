package transform

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/cuda"
)

func GpuScaleUp(uimg cuda.GpuMat, udst *cuda.GpuMat, scale float64) {
	cuda.Resize(uimg, udst, image.Point{0, 0}, scale, scale, cuda.InterpolationDefault)
}

func ScaleUp(img gocv.Mat, dst *gocv.Mat, scale float64) {
	gocv.Resize(img, dst, image.Point{0, 0}, scale, scale, gocv.InterpolationLinear)
}

func GpuAffineTransform(uimg cuda.GpuMat, udst *cuda.GpuMat, w, h int) {
	rot := gocv.GetRotationMatrix2D(image.Point{w / 2, h / 2}, 90.0, 0.75)
	defer rot.Close()
	urot := cuda.NewGpuMatFromMat(rot)
	defer urot.Close()

	cuda.WarpAffine(uimg, udst, urot, image.Point{w, h},
		cuda.InterpolationLinear,
		cuda.BorderDefault,
		color.RGBA{},
	)
}

func AffineTransform(img gocv.Mat, dst *gocv.Mat, w, h int) {
	rot := gocv.GetRotationMatrix2D(image.Point{w / 2, h / 2}, 90.0, 0.75)
	defer rot.Close()
	gocv.WarpAffine(img, dst, rot, image.Point{w, h})
}
