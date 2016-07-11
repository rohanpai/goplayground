
package xxx


import (
	//	&#34;fmt&#34;
	&#34;math&#34;
	&#34;math/rand&#34;

//	&#34;time&#34;
)

type GenFn func(p IntVector)        // Generates a new solution from an old one
type CoolFn func(x float64) float64 // Cooling schedule: Generates a new temperature from the previous one.
//                  Any function that takes a scalar as input and
//                  returns a smaller but positive scalar as output. 

func cool(x float64) float64 {
	return 0.8 * x
}

// SimAnn minimizes a loss function using simulated annealing (Kirkpatrick et al., 1983)
func SimAnn(a Matrix64, loss ObjFn, isLoss bool, proposePerm GenFn, cool CoolFn) (bestEnergy float64, bestSolution IntVector) {
	//    Kirkpatrick, S., Gelatt, C.D., &amp; Vecchi, M.P. (1983). Optimization by Simulated Annealing. Science, 220: 671-680.
	//    Based on &#34;anneal.m&#34; code by Joachim Vandekerckhove joachim.vandekerckhove@psy.kuleuven.be
	//    http://www.mathworks.com/matlabcentral/fileexchange/10548-general-simulated-annealing-algorithm/content/anneal.m

	var (
		seed      int64
		newEnergy float64
	)

	// params
	initT := 1.0 // The initial temperature, can be any positive number.
	minT := 1e-8 // Temperature at which to stop, can be any positive number
	//                  smaller than initT. 
	stopVal := -inf //Value at which to stop immediately, can be any output of
	//                  loss fn that is sufficiently low for you.
	maxConsRej := 1000 // Maximum number of consecutive rejections, can be any
	//                  positive number.
	maxTries := 300 // Maximum number of tries within one temperature, can be
	//                  any positive number.
	maxSuccess := 20 // Maximum number of successes within one temperature, can
	//                  be any positive number.
	k := 1.0 // Boltzmann constant.

	// init
	try := 0
	success := 0
	finished := false
	consec := 0
	temp := initT
	stopVal = -inf
	total := 0
	initialSolution := NewIntVector(a.Rows())
	initialSolution.Perm()
	bestSolution = initialSolution.Clone()
	newSolution := initialSolution.Clone()
	initialEnergy := loss(a, initialSolution)
	bestEnergy = initialEnergy

	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	// seed := time.Now().UnixNano()
	seed = 5
	r := rand.New(rand.NewSource(seed))

	for !finished {
		try&#43;&#43; // just an iteration counter
		newSolution.CopyFrom(bestSolution)

		//  Stop / decrement temp criteria
		if try &gt;= maxTries || success &gt;= maxSuccess {
			if temp &lt; minT || consec &gt;= maxConsRej {
				finished = true
				total = total &#43; try
				break
			} else {
				temp = cool(temp) // decrease temp according to cooling schedule
				total = total &#43; try
				try = 1
				success = 1
			}
		}

		proposePerm(newSolution)
		if isLoss {
			newEnergy = loss(a, newSolution)
		} else {
			newEnergy = -loss(a, newSolution)
		}
		if newEnergy &lt; stopVal {
			bestSolution.CopyFrom(newSolution)
			bestEnergy = newEnergy
			break
		}

		if bestEnergy-newEnergy &gt; 1e-6 {
			bestSolution.CopyFrom(newSolution)
			bestEnergy = newEnergy
			success&#43;&#43;
			consec = 0
		} else {
			if r.Float64() &lt; math.Exp((bestEnergy-newEnergy)/(k*temp)) {
				bestSolution.CopyFrom(newSolution)
				bestEnergy = newEnergy
				success&#43;&#43;
			} else {
				consec&#43;&#43;
			}
		}
	}

	return
}

func RobSA(a Matrix64, objFn ObjFn, isLoss bool) (bestEnergy float64, bestSolution IntVector) {
	bestEnergy, bestSolution = SimAnn(a, objFn, isLoss, proposePerm, cool)
	return
}

// ============= New versions =============

// SimAnn2 minimizes a loss function using simulated annealing (Kirkpatrick et al., 1983)
func SimAnn2(a Matrix64, initialSolution IntVector, loss ObjFn, isLoss bool, proposePerm GenFn, cool CoolFn) (bestEnergy float64, bestSolution IntVector) {
	//    Kirkpatrick, S., Gelatt, C.D., &amp; Vecchi, M.P. (1983). Optimization by Simulated Annealing. Science, 220: 671-680.
	//    Based on &#34;anneal.m&#34; code by Joachim Vandekerckhove joachim.vandekerckhove@psy.kuleuven.be
	//    http://www.mathworks.com/matlabcentral/fileexchange/10548-general-simulated-annealing-algorithm/content/anneal.m

	var (
		seed      int64
		newEnergy float64
	)

	// params
	initT := 1.0 // The initial temperature, can be any positive number.
	minT := 1e-8 // Temperature at which to stop, can be any positive number
	//                  smaller than initT. 
	stopVal := -inf //Value at which to stop immediately, can be any output of
	//                  loss fn that is sufficiently low for you.
	maxConsRej := 1000 // Maximum number of consecutive rejections, can be any
	//                  positive number.
	maxTries := 300 // Maximum number of tries within one temperature, can be
	//                  any positive number.
	maxSuccess := 20 // Maximum number of successes within one temperature, can
	//                  be any positive number.
	k := 1.0 // Boltzmann constant.

	// init
	try := 0
	success := 0
	finished := false
	consec := 0
	temp := initT
	stopVal = -inf
	total := 0
	bestSolution = initialSolution.Clone()
	newSolution := initialSolution.Clone()
	initialEnergy := loss(a, initialSolution)
	bestEnergy = initialEnergy

	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	// seed := time.Now().UnixNano()
	seed = 5
	r := rand.New(rand.NewSource(seed))

	for !finished {
		try&#43;&#43; // just an iteration counter
		newSolution.CopyFrom(bestSolution)

		//  Stop / decrement temp criteria
		if try &gt;= maxTries || success &gt;= maxSuccess {
			if temp &lt; minT || consec &gt;= maxConsRej {
				finished = true
				total = total &#43; try
				break
			} else {
				temp = cool(temp) // decrease temp according to cooling schedule
				total = total &#43; try
				try = 1
				success = 1
			}
		}

		proposePerm(newSolution)
		if isLoss {
			newEnergy = loss(a, newSolution)
		} else {
			newEnergy = -loss(a, newSolution)
		}
		if newEnergy &lt; stopVal {
			bestSolution.CopyFrom(newSolution)
			bestEnergy = newEnergy
			break
		}

		if bestEnergy-newEnergy &gt; 1e-6 {
			bestSolution.CopyFrom(newSolution)
			bestEnergy = newEnergy
			success&#43;&#43;
			consec = 0
		} else {
			if r.Float64() &lt; math.Exp((bestEnergy-newEnergy)/(k*temp)) {
				bestSolution.CopyFrom(newSolution)
				bestEnergy = newEnergy
				success&#43;&#43;
			} else {
				consec&#43;&#43;
			}
		}
	}

	return
}

func RobSA2(a Matrix64, p IntVector, objFn ObjFn, isLoss bool) (bestEnergy float64, bestSolution IntVector) {
	bestEnergy, bestSolution = SimAnn2(a, p, objFn, isLoss, proposePerm, cool)
	return
}

// ==========================================================

// SimAnn3 minimizes a loss function using simulated annealing (Kirkpatrick et al., 1983)
func SimAnn3(a Matrix64, p IntVector, loss ObjFn, isLoss bool, proposePerm GenFn, cool CoolFn) (bestEnergy float64) {
	//    Kirkpatrick, S., Gelatt, C.D., &amp; Vecchi, M.P. (1983). Optimization by Simulated Annealing. Science, 220: 671-680.
	//    Based on &#34;anneal.m&#34; code by Joachim Vandekerckhove joachim.vandekerckhove@psy.kuleuven.be
	//    http://www.mathworks.com/matlabcentral/fileexchange/10548-general-simulated-annealing-algorithm/content/anneal.m

	var (
		seed      int64
		newEnergy float64
	)

	// params
	initT := 1.0 // The initial temperature, can be any positive number.
	minT := 1e-8 // Temperature at which to stop, can be any positive number
	//                  smaller than initT. 
	stopVal := -inf //Value at which to stop immediately, can be any output of
	//                  loss fn that is sufficiently low for you.
	maxConsRej := 1000 // Maximum number of consecutive rejections, can be any
	//                  positive number.
	maxTries := 300 // Maximum number of tries within one temperature, can be
	//                  any positive number.
	maxSuccess := 20 // Maximum number of successes within one temperature, can
	//                  be any positive number.
	k := 1.0 // Boltzmann constant.

	// init
	try := 0
	success := 0
	finished := false
	consec := 0
	temp := initT
	stopVal = -inf
	total := 0
	newSolution := p.Clone()
	initialEnergy := loss(a, p)
	bestEnergy = initialEnergy

	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	// seed := time.Now().UnixNano()
	seed = 5
	r := rand.New(rand.NewSource(seed))

	for !finished {
		try&#43;&#43; // just an iteration counter
		newSolution.CopyFrom(p)

		//  Stop / decrement temp criteria
		if try &gt;= maxTries || success &gt;= maxSuccess {
			if temp &lt; minT || consec &gt;= maxConsRej {
				finished = true
				total = total &#43; try
				break
			} else {
				temp = cool(temp) // decrease temp according to cooling schedule
				total = total &#43; try
				try = 1
				success = 1
			}
		}

		proposePerm(newSolution)
		if isLoss {
			newEnergy = loss(a, newSolution)
		} else {
			newEnergy = -loss(a, newSolution)
		}
		if newEnergy &lt; stopVal {
			p.CopyFrom(newSolution)
			bestEnergy = newEnergy
			break
		}

		if bestEnergy-newEnergy &gt; 1e-6 {
			p.CopyFrom(newSolution)
			bestEnergy = newEnergy
			success&#43;&#43;
			consec = 0
		} else {
			if r.Float64() &lt; math.Exp((bestEnergy-newEnergy)/(k*temp)) {
				p.CopyFrom(newSolution)
				bestEnergy = newEnergy
				success&#43;&#43;
			} else {
				consec&#43;&#43;
			}
		}
	}

	return
}