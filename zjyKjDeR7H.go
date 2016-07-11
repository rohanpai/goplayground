package dist

import (
	&#34;math&#34;

	&#34;github.com/gonum/floats&#34;
	&#34;github.com/gonum/matrix/mat64&#34;
)

// MVNormal represents a multivariate normal distribution with the given mean and
// covariance matrix. Normal must be constructed using NewNormal. Fields of
// MVNormal structure are exported to allow direct use, though they should be treated
// as read-only. If the mean or covariance needs to change, NewNormal must be
// called to generate a new distribution.
type MVNormal struct {
	mu         []float64
	sigma      *mat64.Dense
	logsqrtdet float64
	chol       mat64.CholeskyFactor // Cholesky decomposition of the covariance matrix
	sz         int
}

// NewNormal creates a multivariate normal distribution with mean Mu and
// Covariance matrix Sigma.
func NewNormal(mu []float64, sigma *mat64.Dense) *Normal {
	// TODO (btracey): Replace with a Symmetric matrix when such a type exists
	// and can be used with Cholesky.
	r, c := sigma.Dims()
	if r != c {
		panic(&#34;mvnormal: covariance matrix must be square&#34;)
	}
	if len(m) != r {
		panic(&#34;mvnormal: length of mu must match size of covariance matrix&#34;)
	}

	m := make([]float64, len(mu))
	copy(m, mu)
	s := &amp;mat64.Dense{}
	s.Clone(sigma)
	n := &amp;Normal{
		mu:    mu,
		sigma: sigma,
	}
	n.init()
	return n
}

// init computes some invariants (cholesky decomposition, etc.).
func (n *MVNormal) init() {
	cf := mat64.Cholesky(n.Sigma)

	// Cholesky decomposition doesn&#39;t change the determinant. The determinant
	// of a cholesky matrix is the product of the diagonal values.
	rows, _ := cf.L.Dims()
	var logsqrtdet float64
	for i := 0; i &lt; rows; i&#43;&#43; {
		logsqrtdet &#43;= math.Log(cf.L.At(i, i))
	}
	n.logsqrtdet = 0.5 * logsqrtdet
	n.cholfac = cf
	n.sz = rows
}

func (n *MVNormal) Prob(x []float64) float64 {
	return math.Log(n.LogProb(x))
}

func (n *MVNormal) LogProb(x []float64) float64 {
	dim := len(n.mu)
	c := -0.5*float64(dim)*math.Log(2*math.Pi) - n.logsqrtdet

	// Now need to compute (x-mu)&#39;Sigma^-1 (x-mu)
	xMinusMu := make([]float64, dim)
	floats.SubTo(xMinusMu, x, n.Mu)
	d := mat64.NewDense(dim, 1, xMinusMu)
	// TODO: Replace this with a triangular solve
	d = mat64.Solve(n.cholfac.L, d)
	var sumsq float64
	l := d.RawMatrix().Data
	for i := range l {
		sumsq &#43;= l[i] * l[i]
	}
	return c - 0.5*sumsq
}
