package waves

import "log"

type Config struct {
	Alpha  float64   `json:"alpha"`
	Beta   float64   `json:"beta"`
	Width  float64   `json:"width"`
	Now    []float64 `json:"now"`
	Before []float64 `json:"before,omitempty"`
	Δt     float64   `json:"dt"`
}

// func FromState(f *State, α, β, Δt float64) Config {

// 	cells := make([]float64, f.Domain().Cells())
// 	for i := range cells {
// 		cells[i] = f.At(i)
// 	}

// 	return Config{
// 		Alpha: α, Beta: β, Δt: Δt,
// 		Width: f.Domain().Width(),
// 		Cells: cells,
// 	}
// }

func (cfg Config) InitState() (now, before State) {
	if len(cfg.Now) != len(cfg.Before) {
		log.Fatalf("now has %d cells, while before has %d", len(cfg.Now), len(cfg.Before))
	}
	dom := NewDomain(cfg.Width, len(cfg.Now))
	now = dom.New(func(i int) float64 { return cfg.Now[i] })
	before = dom.New(func(i int) float64 { return cfg.Before[i] })
	return now, before
}
