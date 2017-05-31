package main

func main() {
}

type LeapFrog struct {
	Dom              *Domain
	Next, Curr, Prev State
	F_t              func(i int) float64
	Δx               float64
}

func (leap *LeapFrog) Step() {
	leap.Next = leap.Dom.New(func(i int) float64 {
		return leap.Prev.At(i) + 2*leap.Δx*leap.F_t(i)
	})
	leap.Curr, leap.Prev = leap.Next, leap.Curr
}

func NewDomain(size int) *Domain {
	return &Domain{size}
}

type Domain struct {
	size int
}

type State struct {
	dom  *Domain
	data []float64
}

func (dom *Domain) New(f func(i int) float64) State {
	st := State{dom, make([]float64, dom.size)}
	for i := range st.data {
		st.data[i] = f(i)
	}
	return st
}

func (dom *Domain) Wrap(i int) int {
	for i < 0 {
		i += dom.size
	}
	for i >= dom.size {
		i -= dom.size
	}
	return i
}

func (st *State) At(i int) float64 {
	return st.data[st.dom.Wrap(i)]
}
