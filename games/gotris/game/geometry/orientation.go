package geometry

type (
	Dir  int
	Spin int
)

const (
	DirUnkown Dir = -1
	DirUp     Dir = 0
	DirRight  Dir = 1
	DirDown   Dir = 2
	DirLeft   Dir = 3
)

const (
	SpinUnknown Spin = -1
	SpinLeft    Spin = 0
	SpinRight   Spin = 1
)
