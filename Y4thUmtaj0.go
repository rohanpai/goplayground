// Copyright 2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	&#34;github.com/gonum/matrix/mat64&#34;
	&#34;math&#34;
)

func min(a, b int) int {
	if a &lt; b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a &gt; b {
		return a
	}
	return b
}

// SVD performs singular value decomposition for an m-by-n matrix a with m &gt;= n,
// the singular value decomposition is an m-by-n orthogonal matrix u, an n-by-n
// diagonal matrix s, and an n-by-n orthogonal matrix v so that a = u*s*v&#39;. The
// matrix a is overwritten during the decomposition and u and v are only returned
// when wantu and wantv are true respectively.
//
// The singular values, sigma[k] = s[k][k], are ordered so that
//
//  sigma[0] &gt;= sigma[1] &gt;= ... &gt;= sigma[n-1].
//
// The matrix condition number and the effective numerical rank can be computed from
// this decomposition.
func SVD(a *mat64.Dense, epsilon float64, wantu, wantv bool) (s []float64, u, v *mat64.Dense) {
	m, n := a.Dims()

	// Apparently the failing cases are only a proper subset of (m&lt;n),
	// so let&#39;s not panic. Correct fix to come later?
	// if m &lt; n {
	// 	panic(mat64.ErrShape)
	// }

	s = make([]float64, min(m&#43;1, n))
	nu := min(m, n)
	if wantu {
		u, _ = mat64.NewDense(m, nu, make([]float64, m*nu))
	}
	if wantv {
		v, _ = mat64.NewDense(n, n, make([]float64, n*n))
	}

	var (
		e    = make([]float64, n)
		work = make([]float64, m)
	)

	// Reduce a to bidiagonal form, storing the diagonal elements
	// in s and the super-diagonal elements in e.
	nct := min(m-1, n)
	nrt := max(0, min(n-2, m))
	for k := 0; k &lt; max(nct, nrt); k&#43;&#43; {
		if k &lt; nct {
			// Compute the transformation for the k-th column and
			// place the k-th diagonal in s[k].
			// Compute 2-norm of k-th column without under/overflow.
			s[k] = 0
			for i := k; i &lt; m; i&#43;&#43; {
				s[k] = math.Hypot(s[k], a.At(i, k))
			}
			if s[k] != 0 {
				if a.At(k, k) &lt; 0 {
					s[k] = -s[k]
				}
				for i := k; i &lt; m; i&#43;&#43; {
					a.Set(i, k, a.At(i, k)/s[k])
				}
				a.Set(k, k, a.At(k, k)&#43;1)
			}
			s[k] = -s[k]
		}

		for j := k &#43; 1; j &lt; n; j&#43;&#43; {
			if k &lt; nct &amp;&amp; s[k] != 0 {
				// Apply the transformation.
				var t float64
				for i := k; i &lt; m; i&#43;&#43; {
					t &#43;= a.At(i, k) * a.At(i, j)
				}
				t = -t / a.At(k, k)
				for i := k; i &lt; m; i&#43;&#43; {
					a.Set(i, j, a.At(i, j)&#43;t*a.At(i, k))
				}
			}

			// Place the k-th row of a into e for the
			// subsequent calculation of the row transformation.
			e[j] = a.At(k, j)
		}

		if wantu &amp;&amp; k &lt; nct {
			// Place the transformation in u for subsequent back
			// multiplication.
			for i := k; i &lt; m; i&#43;&#43; {
				u.Set(i, k, a.At(i, k))
			}
		}

		if k &lt; nrt {
			// Compute the k-th row transformation and place the
			// k-th super-diagonal in e[k].
			// Compute 2-norm without under/overflow.
			e[k] = 0
			for i := k &#43; 1; i &lt; n; i&#43;&#43; {
				e[k] = math.Hypot(e[k], e[i])
			}
			if e[k] != 0 {
				if e[k&#43;1] &lt; 0 {
					e[k] = -e[k]
				}
				for i := k &#43; 1; i &lt; n; i&#43;&#43; {
					e[i] /= e[k]
				}
				e[k&#43;1] &#43;= 1.0
			}
			e[k] = -e[k]
			if k&#43;1 &lt; m &amp;&amp; e[k] != 0 {
				// Apply the transformation.
				for i := k &#43; 1; i &lt; m; i&#43;&#43; {
					work[i] = 0
				}
				for j := k &#43; 1; j &lt; n; j&#43;&#43; {
					for i := k &#43; 1; i &lt; m; i&#43;&#43; {
						work[i] &#43;= e[j] * a.At(i, j)
					}
				}
				for j := k &#43; 1; j &lt; n; j&#43;&#43; {
					t := -e[j] / e[k&#43;1]
					for i := k &#43; 1; i &lt; m; i&#43;&#43; {
						a.Set(i, j, a.At(i, j)&#43;t*work[i])
					}
				}
			}
			if wantv {
				// Place the transformation in v for subsequent
				// back multiplication.
				for i := k &#43; 1; i &lt; n; i&#43;&#43; {
					v.Set(i, k, e[i])
				}
			}
		}
	}

	// Set up the final bidiagonal matrix or order p.
	p := min(n, m&#43;1)
	if nct &lt; n {
		s[nct] = a.At(nct, nct)
	}
	if m &lt; p {
		s[p-1] = 0
	}
	if nrt&#43;1 &lt; p {
		e[nrt] = a.At(nrt, p-1)
	}
	e[p-1] = 0

	// If required, generate u.
	if wantu {
		for j := nct; j &lt; nu; j&#43;&#43; {
			for i := 0; i &lt; m; i&#43;&#43; {
				u.Set(i, j, 0)
			}
			u.Set(j, j, 1)
		}
		for k := nct - 1; k &gt;= 0; k-- {
			if s[k] != 0 {
				for j := k &#43; 1; j &lt; nu; j&#43;&#43; {
					var t float64
					for i := k; i &lt; m; i&#43;&#43; {
						t &#43;= u.At(i, k) * u.At(i, j)
					}
					t = -t / u.At(k, k)
					for i := k; i &lt; m; i&#43;&#43; {
						u.Set(i, j, u.At(i, j)&#43;t*u.At(i, k))
					}
				}
				for i := k; i &lt; m; i&#43;&#43; {
					u.Set(i, k, -u.At(i, k))
				}
				u.Set(k, k, 1&#43;u.At(k, k))
				for i := 0; i &lt; k-1; i&#43;&#43; {
					u.Set(i, k, 0)
				}
			} else {
				for i := 0; i &lt; m; i&#43;&#43; {
					u.Set(i, k, 0)
				}
				u.Set(k, k, 1)
			}
		}
	}

	// If required, generate v.
	if wantv {
		for k := n - 1; k &gt;= 0; k-- {
			if k &lt; nrt &amp;&amp; e[k] != 0 {
				for j := k &#43; 1; j &lt; nu; j&#43;&#43; {
					var t float64
					for i := k &#43; 1; i &lt; n; i&#43;&#43; {
						t &#43;= v.At(i, k) * v.At(i, j)
					}
					t = -t / v.At(k&#43;1, k)
					for i := k &#43; 1; i &lt; n; i&#43;&#43; {
						v.Set(i, j, v.At(i, j)&#43;t*v.At(i, k))
					}
				}
			}
			for i := 0; i &lt; n; i&#43;&#43; {
				v.Set(i, k, 0)
			}
			v.Set(k, k, 1)
		}
	}

	// Main iteration loop for the singular values.
	var (
		pp   = p - 1
		iter = 0
		tiny = math.Pow(2, -966.0)
	)
	for p &gt; 0 {
		var k, kase int

		// Here is where a test for too many iterations would go.

		// This section of the program inspects for
		// negligible elements in the s and e arrays.  On
		// completion the variables kase and k are set as follows.
		//
		// kase = 1     if s(p) and e[k-1] are negligible and k&lt;p
		// kase = 2     if s(k) is negligible and k&lt;p
		// kase = 3     if e[k-1] is negligible, k&lt;p, and
		//              s(k), ..., s(p) are not negligible (qr step).
		// kase = 4     if e(p-1) is negligible (convergence).
		//
		for k = p - 2; k &gt;= -1; k-- {
			if k == -1 {
				break
			}
			if math.Abs(e[k]) &lt;= tiny&#43;epsilon*(math.Abs(s[k])&#43;math.Abs(s[k&#43;1])) {
				e[k] = 0
				break
			}
		}
		if k == p-2 {
			kase = 4
		} else {
			var ks int
			for ks = p - 1; ks &gt;= k; ks-- {
				if ks == k {
					break
				}
				var t float64
				if ks != p {
					t = math.Abs(e[ks])
				}
				if ks != k&#43;1 {
					t &#43;= math.Abs(e[ks-1])
				}
				if math.Abs(s[ks]) &lt;= tiny&#43;epsilon*t {
					s[ks] = 0
					break
				}
			}
			if ks == k {
				kase = 3
			} else if ks == p-1 {
				kase = 1
			} else {
				kase = 2
				k = ks
			}
		}
		k&#43;&#43;

		switch kase {
		// Deflate negligible s(p).
		case 1:
			f := e[p-2]
			e[p-2] = 0
			for j := p - 2; j &gt;= k; j-- {
				t := math.Hypot(s[j], f)
				cs := s[j] / t
				sn := f / t
				s[j] = t
				if j != k {
					f = -sn * e[j-1]
					e[j-1] = cs * e[j-1]
				}
				if wantv {
					for i := 0; i &lt; n; i&#43;&#43; {
						t = cs*v.At(i, j) &#43; sn*v.At(i, p-1)
						v.Set(i, p-1, -sn*v.At(i, j)&#43;cs*v.At(i, p-1))
						v.Set(i, j, t)
					}
				}
			}

		// Split at negligible s(k).
		case 2:
			f := e[k-1]
			e[k-1] = 0
			for j := k; j &lt; p; j&#43;&#43; {
				t := math.Hypot(s[j], f)
				cs := s[j] / t
				sn := f / t
				s[j] = t
				f = -sn * e[j]
				e[j] = cs * e[j]
				if wantu {
					for i := 0; i &lt; m; i&#43;&#43; {
						t = cs*u.At(i, j) &#43; sn*u.At(i, k-1)
						u.Set(i, k-1, -sn*u.At(i, j)&#43;cs*u.At(i, k-1))
						u.Set(i, j, t)
					}
				}
			}

		// Perform one qr step.
		case 3:
			// Calculate the shift.
			scale := math.Max(math.Max(math.Max(math.Max(
				math.Abs(s[p-1]), math.Abs(s[p-2])), math.Abs(e[p-2])), math.Abs(s[k])), math.Abs(e[k]),
			)
			sp := s[p-1] / scale
			spm1 := s[p-2] / scale
			epm1 := e[p-2] / scale
			sk := s[k] / scale
			ek := e[k] / scale
			b := ((spm1&#43;sp)*(spm1-sp) &#43; epm1*epm1) / 2.0
			c := (sp * epm1) * (sp * epm1)

			var shift float64
			if b != 0 || c != 0 {
				shift = math.Sqrt(b*b &#43; c)
				if b &lt; 0 {
					shift = -shift
				}
				shift = c / (b &#43; shift)
			}
			f := (sk&#43;sp)*(sk-sp) &#43; shift
			g := sk * ek

			// Chase zeros.
			for j := k; j &lt; p-1; j&#43;&#43; {
				t := math.Hypot(f, g)
				cs := f / t
				sn := g / t
				if j != k {
					e[j-1] = t
				}
				f = cs*s[j] &#43; sn*e[j]
				e[j] = cs*e[j] - sn*s[j]
				g = sn * s[j&#43;1]
				s[j&#43;1] = cs * s[j&#43;1]
				if wantv {
					for i := 0; i &lt; n; i&#43;&#43; {
						t = cs*v.At(i, j) &#43; sn*v.At(i, j&#43;1)
						v.Set(i, j&#43;1, -sn*v.At(i, j)&#43;cs*v.At(i, j&#43;1))
						v.Set(i, j, t)
					}
				}
				t = math.Hypot(f, g)
				cs = f / t
				sn = g / t
				s[j] = t
				f = cs*e[j] &#43; sn*s[j&#43;1]
				s[j&#43;1] = -sn*e[j] &#43; cs*s[j&#43;1]
				g = sn * e[j&#43;1]
				e[j&#43;1] = cs * e[j&#43;1]
				if wantu &amp;&amp; j &lt; m-1 {
					for i := 0; i &lt; m; i&#43;&#43; {
						t = cs*u.At(i, j) &#43; sn*u.At(i, j&#43;1)
						u.Set(i, j&#43;1, -sn*u.At(i, j)&#43;cs*u.At(i, j&#43;1))
						u.Set(i, j, t)
					}
				}
			}
			e[p-2] = f
			iter = iter &#43; 1

		// Convergence.
		case 4:
			// Make the singular values positive.
			if s[k] &lt;= 0 {
				if s[k] &lt; 0 {
					s[k] = -s[k]
				} else {
					s[k] = 0
				}
				if wantv {
					for i := 0; i &lt;= pp; i&#43;&#43; {
						v.Set(i, k, -v.At(i, k))
					}
				}
			}

			// Order the singular values.
			for k &lt; pp {
				if s[k] &gt;= s[k&#43;1] {
					break
				}
				t := s[k]
				s[k] = s[k&#43;1]
				s[k&#43;1] = t
				if wantv &amp;&amp; (k &lt; n-1) {
					for i := 0; i &lt; n; i&#43;&#43; {
						t = v.At(i, k&#43;1)
						v.Set(i, k&#43;1, v.At(i, k))
						v.Set(i, k, t)
					}
				}
				if wantu &amp;&amp; (k &lt; m-1) {
					for i := 0; i &lt; m; i&#43;&#43; {
						t = u.At(i, k&#43;1)
						u.Set(i, k&#43;1, u.At(i, k))
						u.Set(i, k, t)
					}
				}
				k&#43;&#43;
			}
			iter = 0
			p--
		}
	}

	return s, u, v
}

func Rank(a *mat64.Dense, s []float64, epsilon float64) int {
	m, n := a.Dims()
	tol := float64(max(m, n)) * s[0] * epsilon
	r := 0
	for i := 0; i &lt; len(s); i&#43;&#43; {
		if s[i] &gt; tol {
			r&#43;&#43;
		}
	}
	return r
}
