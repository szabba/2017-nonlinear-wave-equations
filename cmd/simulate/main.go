package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/szabba/2017-nonlinear-wave-equations/waves"
)

var (
	L     float64
	cells int
	α, β  float64

	N, dumps int
)

func main() {
	flag.IntVar(&cells, "cells", 50, "the number of cells the domain is divided in")
	flag.IntVar(&cells, "c", 50, "shorthand for -cells")

	flag.Float64Var(&L, "domain-width", 50, "the width of the domain")
	flag.Float64Var(&L, "w", 50, "shorthand for -domain-width")

	flag.Float64Var(&α, "alpha", 0.1, "the α parameter")
	flag.Float64Var(&α, "α", 0.1, "shorthand for -alpha")
	flag.Float64Var(&α, "a", 0.1, "shorthand for -alpha")

	// QUESTION: Should this be computed based on the initial condition instead?
	// ANSWER: No, bc H is not specified!
	flag.Float64Var(&β, "beta", 0.1, "the β parameter")
	flag.Float64Var(&β, "β", 0.1, "shorthand for -beta")
	flag.Float64Var(&β, "b", 0.1, "shorthand for -beta")

	flag.IntVar(&N, "steps", 10000, "number of simulation steps to perform")
	flag.IntVar(&N, "N", 10000, "shorthand for -steps")

	flag.IntVar(&dumps, "dumps", 100, "number of states to dump")
	flag.IntVar(&dumps, "D", 100, "shorthand for -dumps")

	flag.Parse()

	dom := waves.NewDomain(L, cells)

	init, err := ReadState(os.Stdin, dom)
	if err != nil {
		log.Fatal(err)
	}

	var leap LeapFrog
	leap = LeapFrog{
		Dom: dom, Curr: init, Prev: init,
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
	_, err := fmt.Fprintf(w, "dump %d\n", i)
	return err
}

func Dx(f waves.State, i int) float64 {
	dom := f.Domain()
	return 1 / (2 * dom.Δx()) * (f.At(dom.Wrap(i+1)) - f.At(dom.Wrap(i-1)))
}

func D2x(f waves.State, i int) float64 {
	dom := f.Domain()
	return 1 / (2 * dom.Δx()) * (Dx(f, i+1) - Dx(f, i-1))
}

func D3x(f waves.State, i int) float64 {
	dom := f.Domain()
	return 1 / (2 * dom.Δx()) * (D2x(f, i+1) - D2x(f, i-1))
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

func ReadState(r io.Reader, dom *waves.Domain) (waves.State, error) {

	vals := make([]float64, dom.Cells())
	var err error
	for i := 0; i < len(vals) && err == nil; i++ {
		_, err = fmt.Fscan(r, &vals[i])
	}

	if err != nil {
		return waves.State{}, err
	}

	return dom.New(func(i int) float64 { return vals[i] }), nil
}
