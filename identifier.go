package main

import (
	"fmt"

	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/learn"
)

// Identifier represents an identifier network
type Identifier struct {
	n *neural.Network
}

// NewIdentifier returns a *Identifier with given number of inputs
// and a description of its internal structure ([]int)
func NewIdentifier(inputs int, layers []int) *Identifier {
	n := neural.NewNetwork(inputs, layers)
	n.RandomizeSynapses()
	return &Identifier{
		n: n,
	}
}

// Evaluate returns the squared difference
func (i *Identifier) Evaluate(input []byte) float64 {
	inFloat := make([]float64, len(input))
	for i, v := range input {
		inFloat[i] = float64(v) / 255.0
	}
	result := i.n.Calculate(inFloat)[0]
	return result
}

// Learn learns one input ([]byte) according to ideal float64 output
func (i *Identifier) Learn(input []byte, ideal float64) {
	inFloat := make([]float64, len(input))
	for i, v := range input {
		inFloat[i] = float64(v) / 255.0
	}
	learn.Learn(i.n, inFloat, []float64{ideal}, 0.05)
}

// PrintEval prints the evaluation of a []byte input according to ideal float64
func (i *Identifier) PrintEval(input []byte, ideal float64) {
	inFloat := make([]float64, len(input))
	for i, v := range input {
		inFloat[i] = float64(v) / 255.0
	}
	e := learn.Evaluation(i.n, inFloat, []float64{ideal})
	fmt.Println(e)
}
