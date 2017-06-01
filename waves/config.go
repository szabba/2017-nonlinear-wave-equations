package waves

type Config struct {
	Alpha float64   `json:"alpha"`
	Beta  float64   `json:"beta"`
	Width float64   `json:"width"`
	Cells []float64 `json:"cells"`
	Δt    float64   `json:"dt"`
}

func FromState(f *State, α, β, Δt float64) Config {

	cells := make([]float64, f.Domain().Cells())
	for i := range cells {
		cells[i] = f.At(i)
	}

	return Config{
		Alpha: α, Beta: β, Δt: Δt,
		Width: f.Domain().Width(),
		Cells: cells,
	}
}

func (cfg Config) InitState() State {
	return NewDomain(
		cfg.Width, len(cfg.Cells),
	).New(
		func(i int) float64 { return cfg.Cells[i] })
}