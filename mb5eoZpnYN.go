package main

import (
	&#34;fmt&#34;
	&#34;math/big&#34; // high-precision math
)

func arccot(x int64, unity *big.Int) *big.Int {
	bigx := big.NewInt(x)
	xsquared := big.NewInt(x*x)
	sum := big.NewInt(0)
	sum.Div(unity, bigx)
	xpower := big.NewInt(0)
	xpower.Set(sum)
	n := int64(3)
	zero := big.NewInt(0)
	sign := false
	
	term := big.NewInt(0)
	for {
		xpower.Div(xpower, xsquared)
		term.Div(xpower, big.NewInt(n))
		if term.Cmp(zero) == 0 {
			break
		}
		if sign {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}
		sign = !sign
		n &#43;= 2
	}
	return sum
}

func main() {
	ndigits := int64(500)
	digits := big.NewInt(ndigits &#43; 10)
	unity := big.NewInt(0)
	unity.Exp(big.NewInt(10), digits, nil)
	pi := big.NewInt(0)
	four := big.NewInt(4)
	pi.Mul(four, pi.Sub(pi.Mul(four, arccot(5, unity)), arccot(239, unity)))
	//val := big.Mul(4, big.Sub(big.Mul(4, arccot(5, unity)), arccot(239, unity)))
	pistring := pi.String()[0:ndigits]
	fmt.Println(&#34;Computed pi: &#34;, pistring)
	digitcount := make([]int, 10)
	for _, digit := range pistring {
		val := digit - &#39;0&#39;
		digitcount[val]&#43;&#43;
	}
	fmt.Printf(&#34;Digit\tCount\n&#34;)
	for i, digit := range digitcount {
		fmt.Printf(&#34;%d\t%d\n&#34;, i, digit)
	}
}
