package main

import (
	"math/rand"
	"sort"
)

// Gen is interface for generative model
type Gen interface {
	Generate(float64) []byte
	Serialise() []float64
	Deserialise([]float64)
}

// Candidate wraps a Gen and represents it in the population
type Candidate struct {
	g       Gen
	fitness float64
}

// SquaredDistanceFitness updates the fitness of the candidate
// according to its accuracy against supplied func
func (c *Candidate) SquaredDistanceFitness(id *Identifier) {
	var fitness float64
	// idk? 0 to 1
	for i := 0; i < len(alphabet); i++ {
		fitness += id.Evaluate(c.g.Generate(float64(i)))
	}

	c.fitness = fitness / 10
}

// Candidates is a slice of Candidates
// defined for sorting
type Candidates []*Candidate

func (slice Candidates) Len() int {
	return len(slice)
}

func (slice Candidates) Less(i, j int) bool {
	return slice[i].fitness < slice[j].fitness
}

func (slice Candidates) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Evolution is a struct representing the evoltionary process
type Evolution struct {
	population Candidates
	fitness    func(Gen) float64
	newGen     func(Gen, float64) Gen
}

func (e *Evolution) fitnessProportionateSelection() Candidates {

	// get max fitness
	var maxFitness float64
	for _, v := range e.population {
		// fmt.Printf("fitness: %f\n", v.fitness)
		if v.fitness > maxFitness {
			maxFitness = v.fitness
		}
	}

	// fmt.Printf("Max fitness: %f\n", maxFitness)

	var newPop Candidates
	// sort by fitness
	sort.Sort(sort.Reverse(e.population))

	// take the first two by default
	newPop = append(newPop, e.population[:2]...)

	// while there aren't enough new candidates
	for len(newPop) < len(e.population)-5 {
		idx := rand.Intn(len(e.population))
		probability := e.population[idx].fitness / maxFitness
		if rand.Float64() < probability {
			mutatedCandidate := e.newGen(e.population[idx].g, MutationRate)
			newPop = append(newPop, &Candidate{
				g:       mutatedCandidate,
				fitness: 0,
			})
		}
	}

	// add 5 new
	for i := 0; i < 5; i++ {
		newCandidate := e.newGen(e.population[0].g, 1)
		newPop = append(newPop, &Candidate{
			g:       newCandidate,
			fitness: 0,
		})
	}

	return newPop
}

// RunGeneration performs fitness, and reproduction
func (e *Evolution) RunGeneration(fitness *Identifier) {
	// calculate fitness for each candidate
	for _, v := range e.population {
		v.SquaredDistanceFitness(fitness)
	}

	e.population = e.fitnessProportionateSelection()

}

// NewEvolution returns a new evo sim
func NewEvolution(generators []Gen, newGen func(Gen, float64) Gen) *Evolution {
	candidates := make(Candidates, len(generators))
	for i, v := range generators {
		candidates[i] = &Candidate{
			g:       v,
			fitness: 0,
		}
	}
	return &Evolution{
		population: candidates,
		newGen:     newGen,
	}
}
