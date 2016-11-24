// package main
//
// import (
// 	"github.com/NOX73/go-neural"
// 	"github.com/NOX73/go-neural/learn"
// )
//
// // Adversarial is the combination of a generative and identifier network
// // that allows backpropagation
// type Adversarial struct {
// 	identifier *Identifier
// 	generator  *Generator
// 	combined   *neural.Network
// }
//
// // NewAdversarialNetwork returns a new *AdversarialNetwork
// func NewAdversarialNetwork(generator *Generator, identifier *Identifier) *Adversarial {
// 	layers := getLayers(generator, identifier)
// 	network := &neural.Network{
// 		Enters: generator.n.Enters,
// 		Layers: layers,
// 	}
//
// 	network.RandomizeSynapses()
//
// 	// connect generator's out to identifier's in
// 	generator.n.Layers[len(generator.n.Layers)-1].ConnectTo(identifier.n.Layers[0])
//
// 	return &Adversarial{
// 		identifier: identifier,
// 		generator:  generator,
// 		combined:   network,
// 	}
// }
//
// // TrainGenerated invokes the generator, inputs to the identifier on its output
// // as a negative test case; it backpropagates to the whole combined network
// func (a *Adversarial) TrainGenerated(seed float64) {
// 	learn.Learn(a.combined, []float64{seed}, []float64{0}, 0.01)
// }
//
// // TrainReal takes a real input and trains the identifier on it (as a positive
// // case)
// func (a *Adversarial) TrainReal(input []byte) {
// 	a.identifier.Learn(input, 1.0)
// }
//
// func getLayers(generator *Generator, identifier *Identifier) []*neural.Layer {
// 	layers := generator.n.Layers
// 	layers = append(layers, identifier.n.Layers...)
//
// 	return layers
// }

package main

import (
	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/learn"
)

// Adversarial is the combination of a generative and identifier network
// that allows backpropagation
type Adversarial struct {
	identifierEnters []*neural.Enter
	genLayerIndex    int
	n                *neural.Network
}

// NewAdversarialNetwork returns a new *AdversarialNetwork
// args: description of layers (first must be 1), (0-indexed) id of gen layer
func NewAdversarialNetwork(layers []int, genLayerIndex int) *Adversarial {
	n := neural.NewNetwork(1, layers)
	n.RandomizeSynapses()
	n.RandomizeSynapses()
	n.RandomizeSynapses()

	enters := make([]*neural.Enter, layers[genLayerIndex])
	for i := range enters {
		e := neural.NewEnter()
		e.ConnectTo(n.Layers[genLayerIndex])
		enters[i] = e
	}

	return &Adversarial{
		n:                n,
		genLayerIndex:    genLayerIndex,
		identifierEnters: enters,
	}
}

// Identify activates the identification portion of the network
func (a *Adversarial) Identify(enters []float64) []float64 {
	// set & send enters
	for i, e := range a.identifierEnters {
		e.Input = enters[i]
		e.Signal()
	}

	// only activate layers after gen layer
	for _, l := range a.n.Layers[a.genLayerIndex:] {
		l.Calculate()
	}

	outL := a.n.Layers[len(a.n.Layers)-1]
	a.n.Out = make([]float64, len(outL.Neurons))

	for i, neuron := range outL.Neurons {
		a.n.Out[i] = neuron.Out
	}

	return a.n.Out
}

// Generate activates the network and returns the generative slice
func (a *Adversarial) Generate(seed float64) []byte {
	a.n.Calculate([]float64{seed})

	outL := a.n.Layers[a.genLayerIndex]
	outFloats := make([]float64, len(outL.Neurons))

	for i, neuron := range outL.Neurons {
		outFloats[i] = neuron.Out
	}

	outBytes := make([]byte, len(outFloats))
	for i, v := range outFloats {
		outBytes[i] = byte(v * 255)
	}

	return outBytes
}

// TrainGenerated invokes the generator, inputs to the identifier on its output
// as a negative test case; it backpropagates to the whole combined network
func (a *Adversarial) TrainGenerated(seed float64) {
	learn.Learn(a.n, []float64{seed}, []float64{0}, 0.1)
}

// TrainReal takes a real input and trains the identifier on it (as a positive
// case)
func (a *Adversarial) TrainReal(input []byte) {
	inFloat := make([]float64, len(input))
	for i, v := range input {
		inFloat[i] = float64(v) / 255.0
	}
	a.Identify(inFloat)
	a.identifierBackpropagation(a.n, []float64{1}, 0.1)
}

// IdentifierBackpropagation only backpropagates to the gen layer
func (a *Adversarial) identifierBackpropagation(n *neural.Network, ideal []float64, speed float64) {

	idLayers := n.Layers[a.genLayerIndex:]
	deltas := make([][]float64, len(idLayers))

	last := len(idLayers) - 1
	l := idLayers[last]
	deltas[last] = make([]float64, len(l.Neurons))
	for i, n := range l.Neurons {
		deltas[last][i] = n.Out * (1 - n.Out) * (ideal[i] - n.Out)
	}

	for i := last - 1; i >= 0; i-- {
		l := idLayers[i]
		deltas[i] = make([]float64, len(l.Neurons))
		for j, n := range l.Neurons {

			var sum float64
			for k, s := range n.OutSynapses {
				sum += s.Weight * deltas[i+1][k]
			}

			deltas[i][j] = n.Out * (1 - n.Out) * sum
		}
	}

	for i, l := range idLayers {
		for j, n := range l.Neurons {
			for _, s := range n.InSynapses {
				s.Weight += speed * deltas[i][j] * s.In
			}
		}
	}

}
