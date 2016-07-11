// Copyright 2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/gonum/matrix/mat64"
	"math"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SVD performs singular value decomposition for an m-by-n matrix a with m >= n,
// the singular value decomposition is an m-by-n orthogonal matrix u, an n-by-n
// diagonal matrix s, and an n-by-n orthogonal matrix v so that a = u*s*v'. The
// matrix a is overwritten during the decomposition and u and v are only returned
// when wantu and wantv are true respectively.
//
// The singular values, sigma[k] = s[k][k], are ordered so that
//
//  sigma[0] >= sigma[1] >= ... >= sigma[n-1].
//
// The matrix condition number and the effective numerical rank can be computed from
// this decomposition.
func SVD(a *mat64.Dense, epsilon float64, wantu, wantv bool) (s []float64, u, v *mat64.Dense) {
	m, n := a.Dims()

	// Apparently the failing cases are only a proper subset of (m<n),
	// so let's not panic. Correct fix to come later?
	// if m < n {
	// 	panic(mat64.ErrShape)
	// }

	s = make([]float64, min(m+1, n))
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
	for k := 0; k < max(nct, nrt); k++ {
		if k < nct {
			// Compute the transformation for the k-th column and
			// place the k-th diagonal in s[k].
			// Compute 2-norm of k-th column without under/overflow.
			s[k] = 0
			for i := k; i < m; i++ {
				s[k] = math.Hypot(s[k], a.At(i, k))
			}
			if s[k] != 0 {
				if a.At(k, k) < 0 {
					s[k] = -s[k]
				}
				for i := k; i < m; i++ {
					a.Set(i, k, a.At(i, k)/s[k])
				}
				a.Set(k, k, a.At(k, k)+1)
			}
			s[k] = -s[k]
		}

		for j := k + 1; j < n; j++ {
			if k < nct && s[k] != 0 {
				// Apply the transformation.
				var t float64
				for i := k; i < m; i++ {
					t += a.At(i, k) * a.At(i, j)
				}
				t = -t / a.At(k, k)
				for i := k; i < m; i++ {
					a.Set(i, j, a.At(i, j)+t*a.At(i, k))
				}
			}

			// Place the k-th row of a into e for the
			// subsequent calculation of the row transformation.
			e[j] = a.At(k, j)
		}

		if wantu && k < nct {
			// Place the transformation in u for subsequent back
			// multiplication.
			for i := k; i < m; i++ {
				u.Set(i, k, a.At(i, k))
			}
		}

		if k < nrt {
			// Compute the k-th row transformation and place the
			// k-th super-diagonal in e[k].
			// Compute 2-norm without under/overflow.
			e[k] = 0
			for i := k + 1; i < n; i++ {
				e[k] = math.Hypot(e[k], e[i])
			}
			if e[k] != 0 {
				if e[k+1] < 0 {
					e[k] = -e[k]
				}
				for i := k + 1; i < n; i++ {
					e[i] /= e[k]
				}
				e[k+1] += 1.0
			}
			e[k] = -e[k]
			if k+1 < m && e[k] != 0 {
				// Apply the transformation.
				for i := k + 1; i < m; i++ {
					work[i] = 0
				}
				for j := k + 1; j < n; j++ {
					for i := k + 1; i < m; i++ {
						work[i] += e[j] * a.At(i, j)
					}
				}
				for j := k + 1; j < n; j++ {
					t := -e[j] / e[k+1]
					for i := k + 1; i < m; i++ {
						a.Set(i, j, a.At(i, j)+t*work[i])
					}
				}
			}
			if wantv {
				// Place the transformation in v for subsequent
				// back multiplication.
				for i := k + 1; i < n; i++ {
					v.Set(i, k, e[i])
				}
			}
		}
	}

	// Set up the final bidiagonal matrix or order p.
	p := min(n, m+1)
	if nct < n {
		s[nct] = a.At(nct, nct)
	}
	if m < p {
		s[p-1] = 0
	}
	if nrt+1 < p {
		e[nrt] = a.At(nrt, p-1)
	}
	e[p-1] = 0

	// If required, generate u.
	if wantu {
		for j := nct; j < nu; j++ {
			for i := 0; i < m; i++ {
				u.Set(i, j, 0)
			}
			u.Set(j, j, 1)
		}
		for k := nct - 1; k >= 0; k-- {
			if s[k] != 0 {
				for j := k + 1; j < nu; j++ {
					var t float64
					for i := k; i < m; i++ {
						t += u.At(i, k) * u.At(i, j)
					}
					t = -t / u.At(k, k)
					for i := k; i < m; i++ {
						u.Set(i, j, u.At(i, j)+t*u.At(i, k))
					}
				}
				for i := k; i < m; i++ {
					u.Set(i, k, -u.At(i, k))
				}
				u.Set(k, k, 1+u.At(k, k))
				for i := 0; i < k-1; i++ {
					u.Set(i, k, 0)
				}
			} else {
				for i := 0; i < m; i++ {
					u.Set(i, k, 0)
				}
				u.Set(k, k, 1)
			}
		}
	}

	// If required, generate v.
	if wantv {
		for k := n - 1; k >= 0; k-- {
			if k < nrt && e[k] != 0 {
				for j := k + 1; j < nu; j++ {
					var t float64
					for i := k + 1; i < n; i++ {
						t += v.At(i, k) * v.At(i, j)
					}
					t = -t / v.At(k+1, k)
					for i := k + 1; i < n; i++ {
						v.Set(i, j, v.At(i, j)+t*v.At(i, k))
					}
				}
			}
			for i := 0; i < n; i++ {
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
	for p > 0 {
		var k, kase int

		// Here is where a test for too many iterations would go.

		// This section of the program inspects for
		// negligible elements in the s and e arrays.  On
		// completion the variables kase and k are set as follows.
		//
		// kase = 1     if s(p) and e[k-1] are negligible and k<p
		// kase = 2     if s(k) is negligible and k<p
		// kase = 3     if e[k-1] is negligible, k<p, and
		//              s(k), ..., s(p) are not negligible (qr step).
		// kase = 4     if e(p-1) is negligible (convergence).
		//
		for k = p - 2; k >= -1; k-- {
			if k == -1 {
				break
			}
			if math.Abs(e[k]) <= tiny+epsilon*(math.Abs(s[k])+math.Abs(s[k+1])) {
				e[k] = 0
				break
			}
		}
		if k == p-2 {
			kase = 4
		} else {
			var ks int
			for ks = p - 1; ks >= k; ks-- {
				if ks == k {
					break
				}
				var t float64
				if ks != p {
					t = math.Abs(e[ks])
				}
				if ks != k+1 {
					t += math.Abs(e[ks-1])
				}
				if math.Abs(s[ks]) <= tiny+epsilon*t {
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
		k++

		switch kase {
		// Deflate negligible s(p).
		case 1:
			f := e[p-2]
			e[p-2] = 0
			for j := p - 2; j >= k; j-- {
				t := math.Hypot(s[j], f)
				cs := s[j] / t
				sn := f / t
				s[j] = t
				if j != k {
					f = -sn * e[j-1]
					e[j-1] = cs * e[j-1]
				}
				if wantv {
					for i := 0; i < n; i++ {
						t = cs*v.At(i, j) + sn*v.At(i, p-1)
						v.Set(i, p-1, -sn*v.At(i, j)+cs*v.At(i, p-1))
						v.Set(i, j, t)
					}
				}
			}

		// Split at negligible s(k).
		case 2:
			f := e[k-1]
			e[k-1] = 0
			for j := k; j < p; j++ {
				t := math.Hypot(s[j], f)
				cs := s[j] / t
				sn := f / t
				s[j] = t
				f = -sn * e[j]
				e[j] = cs * e[j]
				if wantu {
					for i := 0; i < m; i++ {
						t = cs*u.At(i, j) + sn*u.At(i, k-1)
						u.Set(i, k-1, -sn*u.At(i, j)+cs*u.At(i, k-1))
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
			b := ((spm1+sp)*(spm1-sp) + epm1*epm1) / 2.0
			c := (sp * epm1) * (sp * epm1)

			var shift float64
			if b != 0 || c != 0 {
				shift = math.Sqrt(b*b + c)
				if b < 0 {
					shift = -shift
				}
				shift = c / (b + shift)
			}
			f := (sk+sp)*(sk-sp) + shift
			g := sk * ek

			// Chase zeros.
			for j := k; j < p-1; j++ {
				t := math.Hypot(f, g)
				cs := f / t
				sn := g / t
				if j != k {
					e[j-1] = t
				}
				f = cs*s[j] + sn*e[j]
				e[j] = cs*e[j] - sn*s[j]
				g = sn * s[j+1]
				s[j+1] = cs * s[j+1]
				if wantv {
					for i := 0; i < n; i++ {
						t = cs*v.At(i, j) + sn*v.At(i, j+1)
						v.Set(i, j+1, -sn*v.At(i, j)+cs*v.At(i, j+1))
						v.Set(i, j, t)
					}
				}
				t = math.Hypot(f, g)
				cs = f / t
				sn = g / t
				s[j] = t
				f = cs*e[j] + sn*s[j+1]
				s[j+1] = -sn*e[j] + cs*s[j+1]
				g = sn * e[j+1]
				e[j+1] = cs * e[j+1]
				if wantu && j < m-1 {
					for i := 0; i < m; i++ {
						t = cs*u.At(i, j) + sn*u.At(i, j+1)
						u.Set(i, j+1, -sn*u.At(i, j)+cs*u.At(i, j+1))
						u.Set(i, j, t)
					}
				}
			}
			e[p-2] = f
			iter = iter + 1

		// Convergence.
		case 4:
			// Make the singular values positive.
			if s[k] <= 0 {
				if s[k] < 0 {
					s[k] = -s[k]
				} else {
					s[k] = 0
				}
				if wantv {
					for i := 0; i <= pp; i++ {
						v.Set(i, k, -v.At(i, k))
					}
				}
			}

			// Order the singular values.
			for k < pp {
				if s[k] >= s[k+1] {
					break
				}
				t := s[k]
				s[k] = s[k+1]
				s[k+1] = t
				if wantv && (k < n-1) {
					for i := 0; i < n; i++ {
						t = v.At(i, k+1)
						v.Set(i, k+1, v.At(i, k))
						v.Set(i, k, t)
					}
				}
				if wantu && (k < m-1) {
					for i := 0; i < m; i++ {
						t = u.At(i, k+1)
						u.Set(i, k+1, u.At(i, k))
						u.Set(i, k, t)
					}
				}
				k++
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
	for i := 0; i < len(s); i++ {
		if s[i] > tol {
			r++
		}
	}
	return r
}
