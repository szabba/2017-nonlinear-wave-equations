package main

import (
	"flag"
	"log"
	"math"

	"encoding/json"
	"os"

	"github.com/szabba/2017-nonlinear-wave-equations/waves"
)

func main() {
	var (
		cells int
		cfg   Config
	)

	flag.IntVar(&cells, "cells", 50, "the number of cells the domain is divided in")
	flag.IntVar(&cells, "c", 50, "shorthand for -cells")

	flag.Float64Var(&cfg.Width, "domain-width", 50, "the width of the domain")
	flag.Float64Var(&cfg.Width, "w", 50, "shorthand for -domain-width")

	flag.Float64Var(&cfg.Alpha, "alpha", 0.1, "the α parameter")
	flag.Float64Var(&cfg.Alpha, "α", 0.1, "shorthand for -alpha")
	flag.Float64Var(&cfg.Alpha, "a", 0.1, "shorthand for -alpha")

	// QUESTION: Should this be computed based on the initial condition instead?
	// ANSWER: No, bc H is not specified!
	flag.Float64Var(&cfg.Beta, "beta", 0.1, "the β parameter")
	flag.Float64Var(&cfg.Beta, "β", 0.1, "shorthand for -beta")
	flag.Float64Var(&cfg.Beta, "b", 0.1, "shorthand for -beta")

	flag.Parse()

	cellWidth := cfg.Width / float64(cells)

	cfg.Δt = math.Pow(cellWidth, 3) / 4

	cfg.Now = make([]float64, cells)
	cfg.Before = make([]float64, cells)
	for i := range cfg.Now {
		x0 := -cfg.Width / 2
		xMin := x0 + float64(i)*cellWidth
		xMax := x0 + float64(i+1)*cellWidth
		xMid := (xMax + xMin) / 2

		// FIXME: ?
		// The other program expects a moving ref frame init cond
		// η is expressed in the static ref frame.
		cfg.Now[i] = cfg.η(xMid, 0)
		cfg.Before[i] = cfg.η(xMid, -cfg.Δt)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	err := enc.Encode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	waves.Config
}

func (cfg Config) η(x, t float64) float64 {
	return math.Pow(Sech(math.Sqrt(0.75*cfg.Alpha/cfg.Beta)*(x-(1+cfg.Alpha/2)*t)), 2)
}

func Sech(x float64) float64 {
	return 1 / math.Pow(math.Cosh(x), 2)
}
