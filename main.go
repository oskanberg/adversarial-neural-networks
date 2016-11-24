package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	// Width is the width of the target
	Width = 10
	// Height is the height of the target
	Height = 10
	// Rounds is the number of adversarial rounds to run
	Rounds = 50
	// IdentifierTrainingIterations is the number of iterations to train Identifier
	IdentifierTrainingIterations = 1000
	// NumGenerations is the number of generations to run the sim
	NumGenerations = 500
	// PopulationSize is the number of candidates
	PopulationSize = 20
	// MutationRate is the rate that weights mutate
	MutationRate = 0.1
)

var identifyLayers = []int{Width * Height, 50, 1}
var generateLayers = []int{1, 50, Width * Height}

// var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
// var alphabet = []string{"a", "b", "c", "d", "e", "f"}

var alphabet = []string{"a"}

// GenerateLetters returns data for white-background images with letters
func GenerateLetters() [][]byte {
	var imgs = make([][]byte, len(alphabet))
	for i, char := range alphabet {
		data := UniformWhiteData(Width, Height)
		image := CreateGreyImage(Width, Height, data)
		AddLabel(image, 2, 9, char)
		SaveImage(image, "/Users/oskanberg/Desktop/img/"+char+".png")
		imgs[i] = image.Pix
	}

	return imgs
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// create generators
	generators := make([]Gen, PopulationSize)
	for i := range generators {
		generators[i] = NewGenerator(generateLayers)
	}

	// generator generator
	genGen := func(parent Gen, mutationRate float64) Gen {
		weights := parent.Serialise()
		// mutate
		for i, v := range weights {
			if rand.Float64() < mutationRate {
				// fmt.Printf("previous: %f\n", v)
				delta := rand.Float64() - 0.5
				weights[i] = v + delta
				// fmt.Printf("new: %f\n", weights[i])
			}
		}

		gen := NewGenerator(generateLayers)
		gen.Deserialise(weights)
		return gen
	}

	evo := NewEvolution(generators, genGen)

	data := evo.population[0].g.Generate(0)
	image := CreateGreyImage(Width, Height, data)
	SaveImage(image, "/Users/oskanberg/Desktop/output/initial.png")

	// construct & train identifier
	id := NewIdentifier(Width*Height, identifyLayers)

	fmt.Printf("Generating letters ...\n")
	letters := GenerateLetters()

	fmt.Printf("Training identifier ...\n")
	for i := 0; i < IdentifierTrainingIterations; i++ {
		for _, v := range letters {
			id.Learn(v, 1)
			id.Learn(RandomGreyData(Width, Height), 0)
			id.PrintEval(v, 1)
		}
	}

	for round := 0; round < Rounds; round++ {
		fmt.Printf("Running generations ...\n")
		for i := 0; i < NumGenerations; i++ {
			evo.RunGeneration(id)
		}

		fmt.Printf("Max fitness: %f\n", evo.population[0].fitness)

		fmt.Printf("Training identifier ...\n")
		for i := 0; i < IdentifierTrainingIterations; i++ {
			for i, v := range letters {
				id.Learn(v, 1)
				id.Learn(evo.population[0].g.Generate(float64(i)), 0)
			}
		}
		fmt.Printf("Error revealing fake data:\n")
		id.PrintEval(evo.population[0].g.Generate(0), 0)

		data := evo.population[0].g.Generate(0)
		image := CreateGreyImage(Width, Height, data)
		SaveImage(image, "/Users/oskanberg/Desktop/output/round"+strconv.FormatInt(int64(round), 10)+".png")
	}
}
