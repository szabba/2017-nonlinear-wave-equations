package main

import (
	"flag"
	"io"
	"log"
	"os"

	"encoding/json"

	"github.com/szabba/2017-nonlinear-wave-equations/waves"
)

var (
	N, dumps int

	α, β float64
)

func main() {

	flag.IntVar(&N, "steps", 10000, "number of simulation steps to perform")
	flag.IntVar(&N, "N", 10000, "shorthand for -steps")

	flag.IntVar(&dumps, "dumps", 100, "number of states to dump")
	flag.IntVar(&dumps, "D", 100, "shorthand for -dumps")

	flag.Parse()

	var cfg waves.Config
	err := json.NewDecoder(os.Stdin).Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	α, β = cfg.Alpha, cfg.Beta

	init := cfg.InitState()

	var leap LeapFrog
	leap = LeapFrog{
		Dom: init.Domain(), Curr: init, Prev: init,
		F_t: func(i int) float64 {
			f := leap.Curr
			return 3.0/2.0*α*f.At(i)*Dx(f, i) + 1.0/6.0*β*D3x(f, i)
		},
	}

	Δdump := N / dumps
	for i := 0; i < N; i++ {

		if i%Δdump == 0 {
			err := Dump(os.Stdout, leap.Curr)
			if err != nil {
				log.Fatalf("dump %d failed: %s", i/dumps, err)
			}
		}

		leap.Step()
	}
}

// TODO: Remove
var i = 0

// TODO: Fill in.
func Dump(w io.Writer, f waves.State) error {
	defer func() { i++ }()

	cells := make([]float64, f.Domain().Cells())
	for i := range cells {
		cells[i] = f.At(i)
	}

	cfg := waves.Config{
		Alpha: α, Beta: β,
		Width: f.Domain().Width(),
		Cells: cells,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")

	return enc.Encode(cfg)
}

func Dx(f waves.State, i int) float64 {
	return 1 / (2 * f.Domain().Δx()) * (f.At(i+1) - f.At(i-1))
}

func D2x(f waves.State, i int) float64 {
	return 1 / (2 * f.Domain().Δx()) * (Dx(f, i+1) - Dx(f, i-1))
}

func D3x(f waves.State, i int) float64 {
	return 1 / (2 * f.Domain().Δx()) * (D2x(f, i+1) - D2x(f, i-1))
}

type LeapFrog struct {
	Dom              *waves.Domain
	Next, Curr, Prev waves.State
	F_t              func(i int) float64
	Δt               float64
}

func (leap *LeapFrog) Step() {
	leap.Next = leap.Dom.New(func(i int) float64 {
		return leap.Prev.At(i) + 2*leap.Δt*leap.F_t(i)
	})
	leap.Curr, leap.Prev = leap.Next, leap.Curr
}
