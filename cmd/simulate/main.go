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

	now, before := cfg.InitState()

	var leap LeapFrog
	leap = LeapFrog{
		Δt: cfg.Δt, Dom: now.Domain(),
		Curr: now, Prev: before,

		F_t: func(i int) float64 { return F_t(leap.Curr, i) },
	}

	Δdump := N / dumps
	for i := 0; i < N; i++ {

		if i%Δdump == 0 {
			err := Dump(os.Stdout, leap.Curr, leap.Prev)
			if err != nil {
				log.Fatalf("dump %d failed: %s", i/dumps, err)
			}
		}

		leap.Step()
	}
}

func F_t(f waves.State, i int) float64 {
	return -Dx(f, i) - 3.0/2.0*α*f.At(i)*Dx(f, i) - 1.0/6.0*β*D3x(f, i)
}

// TODO: Remove
var i = 0

// TODO: Fill in.
// * Conserved quantities, incl volume.
// * Derivative
func Dump(w io.Writer, now, before waves.State) error {
	defer func() { i++ }()

	var content struct {
		Now []float64 `json:"now"`
	}

	content.Now = make([]float64, now.Domain().Cells())
	for i := range content.Now {
		content.Now[i] = now.At(i)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")

	return enc.Encode(content)
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
