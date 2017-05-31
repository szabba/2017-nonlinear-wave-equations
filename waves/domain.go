package waves

type Domain struct {
	width float64
	cells int
}

func NewDomain(width float64, cells int) *Domain {
	return &Domain{width, cells}
}

func (dom *Domain) Δx() float64 { return dom.Width() / float64(dom.Cells()) }

func (dom *Domain) Width() float64 { return dom.width }

func (dom *Domain) Cells() int { return dom.cells }

func (dom *Domain) Wrap(i int) int {
	for i < 0 {
		i += dom.cells
	}
	for i >= dom.cells {
		i -= dom.cells
	}
	return i
}
