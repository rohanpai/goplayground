package main

import &#34;fmt&#34;
import . &#34;math&#34;
import StlRand &#34;math/rand&#34;
//import &#34;math/Rand&#34;

// Set of functions for the normal distribution defined as
// p(x) = 1/(σ √(2 π)) exp( - (x - μ)² / σ²)

const (
	log2pi = 1.83787706640934548356065947281123527972279494727556682563430308096553139185452079538948659727190839524
)

type Norm struct {
	// Normal distribution with a fixed mean and variance. 
	// Create using New(mean,std)
	μ	float64
	σ	float64
}

func New(μ, σ float64) *Norm {
	var n = new(Norm)
	n.SetParams([]float64{μ, σ})
	return n
}

func (n *Norm) SetParams(p []float64) {
	n.SetMean(p[0])
	n.SetStd(p[1])
	return
}

func (norm *Norm) SetMean(μ float64) {
	// A getter function for the mean of the distribution
	norm.μ = μ
	return
}

func (norm *Norm) SetStd(σ float64) {
	// A getter function for the standard deviation of the distribution

	// Not sure what to do about σ &lt; 0. paanic? pass error?
	norm.σ = σ
	return
}

func (norm *Norm) Rand() float64 {
	// Generates a random variable
	return Rand(norm.μ, norm.σ)
}

func (norm *Norm) Lprob(x float64) float64 {
	// Computes the log of the probability of x
	return Lprob(x, norm.μ, norm.σ)
}

func Rand(μ, σ float64) float64 {
	// Generates a normal random variable which has mean μ and
	// standard deviation σ
	return StlRand.NormFloat64()*σ &#43; μ
}

func Lprob(x, μ, σ float64) float64 {
	// Computes the log of the probability of x for the normal
	// distribution with mean μ and standard deviation σ
	return -Log(σ) - 0.5*log2pi - Pow(((x-μ)/σ), 2)/2.0
}

func main() {
	fmt.Println(&#34;Hello, playground&#34;)
	v := New(5.0, 6.0)
	x := v.Rand()
	fmt.Println(v.Lprob(x))
}
