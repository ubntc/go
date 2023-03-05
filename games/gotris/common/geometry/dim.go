package geometry

type Dim struct {
	W int
	H int
}

func NewDim(w, h int) *Dim { return &Dim{w, h} }
