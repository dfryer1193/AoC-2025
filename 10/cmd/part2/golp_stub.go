//go:build !golp
// +build !golp

package main

// solveGolp is a stub for when the 'golp' build tag is not used.
func solveGolp(masks []int, target []int) (int, bool) {
	return 0, false
}