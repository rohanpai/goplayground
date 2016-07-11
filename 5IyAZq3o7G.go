// Set of functions for the normal distribution defined as
// p(x) = 1/(σ √(2 π)) exp( - (x - μ)² / σ²)

const (
	log2pi = 1.83787706640934548356065947281123527972279494727556682563430308096553139185452079538948659727190839524
)

func z(x,μ, σ float64) float64{
	return (x - μ)/Sqrt(2*Pow(σ,2))
}

func Rand(μ, σ float64) float64{
	// Generates a normal random variable which has mean μ and
	// standard deviation σ
	return StlRand.NormFloat64() * σ &#43; μ
}

func Lprob(x, μ, σ float64) float64{
	// Computes the log of the probability of x for the normal
	// distribution with mean μ and standard deviation σ
	return - Log(σ) - 0.5 * log2pi - Pow(((x - μ)/σ),2)/2.0
}

func Prob(x, μ, σ float64) float64{
	// Computes the probability of x for the normal
	// distribution with mean μ and standard deviation σ
	return Exp(Lprob(x,μ,σ))
}

func Cdf(x, μ, σ float64) float64{
	// Computes the cdf
	return 0.5 &#43; 0.5 * Erf(z(x, μ, σ))
}

func Lcdf(x, μ, σ float64) float64{
	// Computes the log of the cdf
	// THIS NEEDS TO BE CHANGED TO WORK WITH SMALL Z
	return Log(Cdf(x, μ, σ))
}

func Ccdf(x, μ, σ float64) float64{
	// Computes the complementary cdf
	return 0.5 * Erfc(z(x, μ, σ))
}

func Lccdf(x, μ, σ float64) float64{
	// THIS NEEDS TO BE CHANGED TO WORK WITH LARGE Z
	return Log(Ccdf(x, μ, σ))
}

func Icdf(x, μ, σ float64) float64{
	// Computes the inverse cdf
	// NEED TO CODE
	return 1.0
}



// ppf -- inverse cdf (percent point function)
// isf -- inverse survival function
// moment (non-central moment of order n)
// entropy (differential) entropy of the rv
// fit -- parameter estimates for generic data
// mean
// median 
// mode
// variance
// standard deviation
// interval -- not sure what that is
// Derivatives of all of them with respect to the parameters and x

type Norm struct{
	// Normal distribution with a fixed mean and variance. 
	// Create using New(mean,std)
	μ float64
	σ float64
}

func New(μ, σ float64) *Norm{
	var n = new(Norm)
	n.SetParams([]float64{μ, σ})
	return n
}

func (norm *Norm) Mean() float64{
	// A getter function for the mean of the distribution
	return norm.μ
}

func (norm *Norm) Std() float64{
	// A getter function for the standard deviation of the distribution
	return norm.σ
}

func (norm *Norm) Params() []float64{
	return []float64{norm.μ, norm.σ}
}

func (norm *Norm) SetMean(μ float64){
	// A getter function for the mean of the distribution
	norm.μ = μ
	return
}

func (norm *Norm) SetStd(σ float64){
	// A getter function for the standard deviation of the distribution
	
	// Not sure what to do about σ &lt; 0. paanic? pass error?
	norm.σ = σ
	return
}

func (n *Norm) SetParams(p []float64){
	n.SetMean(p[0])
	n.SetStd(p[1])
	return
}

func (norm *Norm)Rand() float64{
	// Generates a random variable
	return Rand(norm.μ, norm.σ)
}

func (norm *Norm)Lprob(x float64) float64{
	// Computes the log of the probability of x
	return Lprob(x, norm.μ, norm.σ)
}

func (norm *Norm)Prob(x float64) float64{
	// Computes the probability of x
	return Lprob(x, norm.μ, norm.σ)
}

func (norm *Norm)Cdf(x float64) float64{
	// Computes the probability of x
	return Cdf(x, norm.μ, norm.σ)
}
func (norm *Norm)Logcdf(x float64) float64{
	// Computes the probability of x
	return Lcdf(x, norm.μ, norm.σ)
}