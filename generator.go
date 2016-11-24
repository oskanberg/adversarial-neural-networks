package main

import (
	"fmt"

	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/learn"
)

// Generator is a generative network
type Generator struct {
	n *neural.Network
}

// NewGenerator returns a new *Generator
func NewGenerator(layers []int) *Generator {
	n := neural.NewNetwork(1, layers)
	n.RandomizeSynapses()
	return &Generator{
		n: n,
	}
}

// Learn takes an id (float64) and a slice of ideal outputs ([]byte)
func (g *Generator) Learn(id float64, ideal []byte) {
	idealOut := make([]float64, len(ideal))
	for i, v := range ideal {
		idealOut[i] = float64(v) / 255.0
	}
	learn.Learn(g.n, []float64{id}, idealOut, 0.01)
}

// PrintEval prints the result of an evaluation on a given id
// and against ideal byte slice
func (g *Generator) PrintEval(id float64, ideal []byte) {
	idealOut := make([]float64, len(ideal))
	for i, v := range ideal {
		idealOut[i] = float64(v) / 255.0
	}
	e := learn.Evaluation(g.n, []float64{id}, idealOut)
	fmt.Println(e)
}

// Generate returns the generated byte slice for a given id
func (g *Generator) Generate(id float64) []byte {
	outputFloats := g.n.Calculate([]float64{id})

	outBytes := make([]byte, len(outputFloats))
	for i, v := range outputFloats {
		outBytes[i] = byte(v * 255)
	}

	return outBytes
}

// Serialise exports the weights of the generator
func (g *Generator) Serialise() []float64 {
	var weights []float64
	for _, l := range g.n.Layers {
		for _, n := range l.Neurons {
			for _, s := range n.InSynapses {
				weights = append(weights, s.Weight)
			}
		}
	}
	return weights
}

// Deserialise updates the generator with given weights
func (g *Generator) Deserialise(weights []float64) {
	var i int
	for _, l := range g.n.Layers {
		for _, n := range l.Neurons {
			for _, s := range n.InSynapses {
				s.Weight = weights[i]
				i++
			}
		}
	}
}
