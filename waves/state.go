package waves

type State struct {
	dom  *Domain
	data []float64
}

func (dom *Domain) New(f func(i int) float64) State {
	st := State{dom, make([]float64, dom.cells)}
	for i := range st.data {
		st.data[i] = f(i)
	}
	return st
}

func (st State) Domain() *Domain { return st.dom }

func (st State) At(i int) float64 {
	return st.data[st.dom.Wrap(i)]
}

func (st State) ToSlice() []float64 {
	out := make([]float64, st.Domain().Cells())
	copy(out, st.data)
	return out
}
