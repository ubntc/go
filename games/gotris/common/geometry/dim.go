package geometry

type Dim struct {
	W int `json:"w,omitempty"`
	H int `json:"h,omitempty"`
}

func NewDim(w, h int) *Dim { return &Dim{w, h} }
