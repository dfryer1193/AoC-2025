//go:build golp
// +build golp

package main

import (
	"math"
	"github.com/draffensperger/golp"
)

// solveGolp uses the GLPK library via the golp package to solve the integer program.
func solveGolp(masks []int, target []int) (int, bool) {
	d := len(target)
	n := len(masks)

	// Quick checks
	if n == 0 {
		for _, v := range target {
			if v != 0 {
				return 0, false
			}
		}
		return 0, true
	}

	lp := golp.NewLP(d, n)

	// Set Objective: Minimize sum(x_j)
	objFn := make([]float64, n)
	for j := 0; j < n; j++ {
		objFn[j] = 1.0
	}
	lp.SetObjFn(objFn)
	// Default is minimize, so no need to call SetMaximize()

	// Define Variables (Columns): non-negative integers
	for j := 0; j < n; j++ {
		lp.SetInt(j, true) // Make variable integer
		lp.SetBounds(j, 0.0, math.MaxFloat64) // Non-negative, effectively unbounded upper
	}

	// Add Constraints (Rows): sum(A_ij * x_j) = t_i
	for i := 0; i < d; i++ {
		rowCoeffs := make([]float64, n)
		for j := 0; j < n; j++ {
			if ((masks[j] >> uint(i)) & 1) == 1 {
				rowCoeffs[j] = 1.0
			} else {
				rowCoeffs[j] = 0.0
			}
		}
		lp.AddConstraint(rowCoeffs, golp.EQ, float64(target[i]))
	}

	// Solve
	lp.SetVerboseLevel(golp.NEUTRAL) // Suppress verbose output
	status := lp.Solve()

	if status == golp.OPTIMAL {
		return int(lp.Objective() + 0.5), true // Round to nearest integer
	}

	return 0, false
}
