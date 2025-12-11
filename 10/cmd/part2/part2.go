package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type machine struct {
	reqMask  int
	buttons  []int
	joltages []int
}

// Try to find the shortest sequence of button presses that configures the machine.
// If joltages are present, each button increments listed counters by 1 and counters start at 0.
// Otherwise, fall back to the light-toggle XOR model.
func (m *machine) Start() int {
	// Joltages mode: counters with non-negative integer presses; use DFS on button press counts with pruning.
	if len(m.joltages) > 0 {
		target := m.joltages
		d := len(target)

		// ensure every required counter can be influenced by some button
		covered := make([]bool, d)
		for _, btn := range m.buttons {
			for i := 0; i < d; i++ {
				if (btn>>uint(i))&1 == 1 {
					covered[i] = true
				}
			}
		}
		for i := 0; i < d; i++ {
			if target[i] > 0 && !covered[i] {
				return -1
			}
		}

		// dedupe and sort buttons by descending bitcount to improve pruning
		maskSeen := map[int]bool{}
		masks := []int{}
		for _, b := range m.buttons {
			if !maskSeen[b] {
				maskSeen[b] = true
				masks = append(masks, b)
			}
		}
		sort.Slice(masks, func(i, j int) bool {
			return bits.OnesCount(uint(masks[i])) > bits.OnesCount(uint(masks[j]))
		})

		INF := int(1e9)
		rem := make([]int, d)
		copy(rem, target)

		// compute multipliers to encode rem -> key
		mult := make([]int, d)
		base := 1
		for i := d - 1; i >= 0; i-- {
			mult[i] = base
			base *= (target[i] + 1)
		}

		// precompute suffix max bits to help LB pruning
		suffixMaxBits := make([]int, len(masks)+1)
		for i := len(masks) - 1; i >= 0; i-- {
			cnt := bits.OnesCount(uint(masks[i]))
			suffixMaxBits[i] = suffixMaxBits[i+1]
			if cnt > suffixMaxBits[i] {
				suffixMaxBits[i] = cnt
			}
		}

		memo := map[uint64]int{}
		var solve func(idx int, rem []int) int
		solve = func(idx int, rem []int) int {
			// encode rem
			keyInt := 0
			sumRem := 0
			maxRem := 0
			for i := 0; i < d; i++ {
				v := rem[i]
				keyInt += v * mult[i]
				sumRem += v
				if v > maxRem {
					maxRem = v
				}
			}
			if keyInt == 0 {
				return 0
			}
			if idx >= len(masks) {
				return INF
			}
			combined := (uint64(keyInt) << 8) | uint64(idx)
			if v, ok := memo[combined]; ok {
				return v
			}
			// lower bound: need at least maxRem presses (one per counter), and at least ceil(sumRem / maxBits)
			lb := maxRem
			maxBits := suffixMaxBits[idx]
			if maxBits > 0 {
				alt := (sumRem + maxBits - 1) / maxBits
				if alt > lb {
					lb = alt
				}
			}

			bestLocal := INF
			btn := masks[idx]
			hasBit := false
			maxT := INF
			for i := 0; i < d; i++ {
				if (btn>>uint(i))&1 == 1 {
					hasBit = true
					if rem[i] < maxT {
						maxT = rem[i]
					}
				}
			}
			if !hasBit {
				res := solve(idx+1, rem)
				memo[combined] = res
				return res
			}
			for t := maxT; t >= 0; t-- {
				// quick pruning
				if t+lb >= bestLocal {
					continue
				}
				if t > 0 {
					for i := 0; i < d; i++ {
						if (btn>>uint(i))&1 == 1 {
							rem[i] -= t
						}
					}
				}
				sub := solve(idx+1, rem)
				if sub != INF {
					cand := sub + t
					if cand < bestLocal {
						bestLocal = cand
					}
				}
				if t > 0 {
					for i := 0; i < d; i++ {
						if (btn>>uint(i))&1 == 1 {
							rem[i] += t
						}
					}
				}
				// early exit if we reached LB
				if bestLocal == lb {
					break
				}
			}
			memo[combined] = bestLocal
			return bestLocal
		}

		// Try GLPK via golp package
		if v, ok := solveGolp(masks, target); ok {
			return v
		}
		// Try partition-DFS solver inspired by Reddit (choose counter with fewest buttons, iterate partitions)
		if v, ok := solvePartitionDFS(masks, target); ok {
			return v
		}
		// Try branch-and-bound native solver first
		if v, ok := solveBnB(masks, target); ok {
			return v
		}
		// fallback to native solver
		res := solve(0, rem)
		if res == INF {
			return -1
		}
		return res
	}

	// Lights/toggle mode (bitmask BFS)
	if m.reqMask == 0 {
		return 0
	}

	level := []int{0}
	visited := map[int]bool{0: true}
	presses := 0

	for len(level) > 0 {
		presses++
		nextLevel := []int{}
		for _, mask := range level {
			for _, button := range m.buttons {
				nextMask := mask ^ button
				if nextMask == m.reqMask {
					return presses
				}
				if !visited[nextMask] {
					visited[nextMask] = true
					nextLevel = append(nextLevel, nextMask)
				}
			}
		}
		level = nextLevel
	}
	return -1
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: go run part1.go <filename>")
		return
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	machines := make([]*machine, 0)
	minPresses := 0

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		machine := &machine{}

		if len(parts) > 0 {
			machine.reqMask = parseLights(parts[0])
		}

		for i := 1; i < len(parts); i++ {
			part := parts[i]
			if strings.HasPrefix(part, "(") {
				machine.buttons = append(machine.buttons, parseButtons(part))
			} else if strings.HasPrefix(part, "{") {
				machine.joltages = parseJoltages(part)
			}
		}

		machines = append(machines, machine)
		minPresses += machine.Start()
	}

	fmt.Println("Total minimum button presses for all machines:", minPresses)
}

func parseLights(state string) int {
	lights := 0
	lightStr := state[1 : len(state)-1]
	for i, ch := range lightStr {
		if ch == '#' {
			lights |= (1 << uint(i))
		}
	}
	return lights
}

func parseButtons(state string) int {
	buttonMask := 0
	nums := strings.Split(state[1:len(state)-1], ",")
	for _, num := range nums {
		val, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("Error parsing button value:", err)
			continue
		}
		buttonMask |= (1 << val)
	}
	return buttonMask
}

func parseJoltages(state string) []int {
	joltages := make([]int, 0)
	nums := strings.Split(state[1:len(state)-1], ",")
	for _, num := range nums {
		val, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("Error parsing joltage value:", err)
			continue
		}
		joltages = append(joltages, val)
	}
	return joltages
}

// solveBnB implements a branch-and-bound native solver for the integer system A x = target.
// Masks are button bitmasks (length n), target is length d. Returns (value, true) if solved.
func solveBnB(masks []int, target []int) (int, bool) {
	d := len(target)
	n := len(masks)
	INF := int(1e9)
	if n == 0 {
		// if target all zeros, zero presses, otherwise impossible
		for _, v := range target {
			if v != 0 {
				return 0, false
			}
		}
		return 0, true
	}
	// compute multipliers for encoding rem -> int key
	mult := make([]int, d)
	base := 1
	for i := d - 1; i >= 0; i-- {
		mult[i] = base
		base *= (target[i] + 1)
	}

	start := time.Now()
	limit := 30 * time.Second
	// greedy upper bound
	rem := make([]int, d)
	copy(rem, target)
	sumRem := func(arr []int) int {
		s := 0
		for _, v := range arr {
			s += v
		}
		return s
	}
	ub := 0
	for sumRem(rem) > 0 {
		if time.Since(start) > limit {
			return 0, false
		}
		bestIdx := -1
		bestCover := 0
		for j, btn := range masks {
			cover := 0
			for i := 0; i < d; i++ {
				if rem[i] > 0 && ((btn>>uint(i))&1) == 1 {
					cover++
				}
			}
			if cover > bestCover {
				bestCover = cover
				bestIdx = j
			}
		}
		if bestCover == 0 {
			return 0, false
		}
		// press it enough to zero one of its covered counters
		minRem := INF
		for i := 0; i < d; i++ {
			if ((masks[bestIdx] >> uint(i)) & 1) == 1 {
				if rem[i] < minRem {
					minRem = rem[i]
				}
			}
		}
		if minRem == INF {
			return 0, false
		}
		for i := 0; i < d; i++ {
			if ((masks[bestIdx] >> uint(i)) & 1) == 1 {
				if rem[i] >= minRem {
					rem[i] -= minRem
				}
			}
		}
		ub += minRem
	}

	// memoization map: key = (encoded rem << 8) | idx
	memo := map[uint64]int{}

	// precompute suffixMaxBits to compute LB quickly
	suffixMaxBits := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		c := bits.OnesCount(uint(masks[i]))
		suffixMaxBits[i] = suffixMaxBits[i+1]
		if c > suffixMaxBits[i] {
			suffixMaxBits[i] = c
		}
	}

	var dfs func(idx int, rem []int) int
	dfs = func(idx int, rem []int) int {
		if time.Since(start) > limit {
			return INF
		}
		// encode rem
		keyInt := 0
		sum := 0
		maxRem := 0
		for i := 0; i < d; i++ {
			v := rem[i]
			keyInt += v * mult[i]
			sum += v
			if v > maxRem {
				maxRem = v
			}
		}
		if keyInt == 0 {
			return 0
		}
		if idx >= n {
			return INF
		}
		combined := (uint64(keyInt) << 8) | uint64(idx)
		if v, ok := memo[combined]; ok {
			return v
		}
		// lower bound
		lb := maxRem
		maxBits := suffixMaxBits[idx]
		if maxBits > 0 {
			alt := (sum + maxBits - 1) / maxBits
			if alt > lb {
				lb = alt
			}
		}
		if lb >= ub {
			memo[combined] = INF
			return INF
		}

		bestLocal := INF
		btn := masks[idx]
		hasBit := false
		maxT := INF
		for i := 0; i < d; i++ {
			if (btn>>uint(i))&1 == 1 {
				hasBit = true
				if rem[i] < maxT {
					maxT = rem[i]
				}
			}
		}
		if !hasBit {
			res := dfs(idx+1, rem)
			memo[combined] = res
			return res
		}
		// try t from maxT down to 0
		for t := maxT; t >= 0; t-- {
			// simple pruning
			if t+lb >= bestLocal || t >= ub {
				continue
			}
			if t > 0 {
				for i := 0; i < d; i++ {
					if (btn>>uint(i))&1 == 1 {
						rem[i] -= t
					}
				}
			}
			sub := dfs(idx+1, rem)
			if sub != INF {
				cand := sub + t
				if cand < bestLocal {
					bestLocal = cand
				}
			}
			if t > 0 {
				for i := 0; i < d; i++ {
					if (btn>>uint(i))&1 == 1 {
						rem[i] += t
					}
				}
			}
			if bestLocal == lb {
				break
			}
		}
		memo[combined] = bestLocal
		return bestLocal
	}

	res := dfs(0, target)
	if res >= INF {
		return 0, false
	}
	return res, true
}

// solvePartitionDFS implements the Reddit approach: pick the counter affected by fewest buttons
// (tie-breaker: largest remaining), enumerate partitions of that counter's value across its buttons,
// apply assignments (pruning if any counter goes negative), recurse with memoization.
func solvePartitionDFS(masks []int, target []int) (int, bool) {
	d := len(target)
	n := len(masks)
	INF := int(1e9)
	// quick checks
	if n == 0 {
		for _, v := range target {
			if v != 0 {
				return 0, false
			}
		}
		return 0, true
	}

	// precompute which buttons affect which counters
	btnsFor := make([][]int, d)
	for j := 0; j < n; j++ {
		for i := 0; i < d; i++ {
			if (masks[j]>>uint(i))&1 == 1 {
				btnsFor[i] = append(btnsFor[i], j)
			}
		}
	}

	// multipliers for encoding rem into an integer key
	mult := make([]int, d)
	base := 1
	for i := d - 1; i >= 0; i-- {
		mult[i] = base
		base *= (target[i] + 1)
	}

	memo := map[int]int{}
	// timeout to avoid excessive runtimes
	start := time.Now()
	limit := 28 * time.Second

	var dfs func(rem []int) int
	dfs = func(rem []int) int {
		if time.Since(start) > limit {
			return INF
		}
		// encode
		key := 0
		have := false
		for i := 0; i < d; i++ {
			key += rem[i] * mult[i]
			if rem[i] != 0 {
				have = true
			}
		}
		if !have {
			return 0
		}
		if v, ok := memo[key]; ok {
			return v
		}

		// choose counter with smallest number of affecting buttons (and largest rem if tie)
		bestI := -1
		bestCnt := 1 << 30
		bestRem := -1
		for i := 0; i < d; i++ {
			if rem[i] == 0 {
				continue
			}
			cnt := len(btnsFor[i])
			if cnt == 0 {
				memo[key] = INF
				return INF
			}
			if cnt < bestCnt || (cnt == bestCnt && rem[i] > bestRem) {
				bestCnt = cnt
				bestI = i
				bestRem = rem[i]
			}
		}

		buttons := btnsFor[bestI]
		k := len(buttons)
		val := rem[bestI]

		bestLocal := INF

		// partial assignment array
		assigned := make([]int, k)

		// helper to check partial assigned doesn't exceed any rem for any counter
		checkPartial := func() bool {
			for c := 0; c < d; c++ {
				sum := 0
				for idx := 0; idx < k; idx++ {
					if ((masks[buttons[idx]] >> uint(c)) & 1) == 1 {
						sum += assigned[idx]
						if sum > rem[c] {
							return false
						}
					}
				}
			}
			return true
		}

		var gen func(pos int, left int)
		gen = func(pos int, left int) {
			if pos == k {
				if left != 0 {
					return
				}
				// apply assigned to create new rem
				newRem := make([]int, d)
				copy(newRem, rem)
				for idx := 0; idx < k; idx++ {
					t := assigned[idx]
					if t == 0 {
						continue
					}
					btn := masks[buttons[idx]]
					for c := 0; c < d; c++ {
						if ((btn >> uint(c)) & 1) == 1 {
							newRem[c] -= t
							if newRem[c] < 0 {
								return
							}
						}
					}
				}
				// recurse
				sub := dfs(newRem)
				if sub != INF {
					sumAssigned := 0
					for _, v := range assigned {
						sumAssigned += v
					}
					if sub+sumAssigned < bestLocal {
						bestLocal = sub + sumAssigned
					}
				}
				return
			}
			// pruning: simple check
			for t := left; t >= 0; t-- {
				assigned[pos] = t
				if !checkPartial() {
					continue
				}
				gen(pos+1, left-t)
				assigned[pos] = 0
				// early exit if bestLocal is zero
				if bestLocal == 0 {
					return
				}
			}
		}

		gen(0, val)
		memo[key] = bestLocal
		return bestLocal
	}

	res := dfs(target)
	if res >= INF {
		return 0, false
	}
	return res, true
}

// solveZ3 uses the z3 binary to solve the integer system A x = target minimizing sum(x).
func solveZ3(masks []int, target []int) (int, bool) {
	n := len(masks)
	d := len(target)
	// quick checks
	if n == 0 {
		for _, v := range target {
			if v != 0 {
				return 0, false
			}
		}
		return 0, true
	}
	var b strings.Builder
	for j := 0; j < n; j++ {
		b.WriteString(fmt.Sprintf("(declare-const x%d Int)\n", j))
		b.WriteString(fmt.Sprintf("(assert (>= x%d 0))\n", j))
	}
	for i := 0; i < d; i++ {
		terms := []string{}
		for j := 0; j < n; j++ {
			if ((masks[j] >> uint(i)) & 1) == 1 {
				terms = append(terms, fmt.Sprintf("x%d", j))
			}
		}
		if len(terms) == 0 {
			if target[i] != 0 {
				return 0, false
			}
			continue
		}
		// build sum term
		if len(terms) == 1 {
			b.WriteString(fmt.Sprintf("(assert (= %s %d))\n", terms[0], target[i]))
		} else {
			b.WriteString("(assert (= (+ ")
			for idx, t := range terms {
				if idx > 0 {
					b.WriteString(" ")
				}
				b.WriteString(t)
			}
			b.WriteString(fmt.Sprintf(") %d))\n", target[i]))
		}
	}
	// objective
	vars := []string{}
	for j := 0; j < n; j++ {
		vars = append(vars, fmt.Sprintf("x%d", j))
	}
	b.WriteString(fmt.Sprintf("(minimize (+ %s))\n", strings.Join(vars, " ")))
	b.WriteString("(check-sat)\n(get-model)\n")

	cmd := exec.Command("z3", "-in")
	cmd.Stdin = strings.NewReader(b.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error running z3:", err)
		return 0, false
	}
	s := string(out)
	if !strings.Contains(s, "sat") {
		return 0, false
	}
	re := regexp.MustCompile(`\(define-fun\s+(x\d+)\s+\(\)\s+Int\s+(-?\d+)\)`)
	matches := re.FindAllStringSubmatch(s, -1)
	vals := map[string]int{}
	for _, m := range matches {
		v, err := strconv.Atoi(m[2])
		if err != nil {
			return 0, false
		}
		vals[m[1]] = v
	}
	total := 0
	for j := 0; j < n; j++ {
		name := fmt.Sprintf("x%d", j)
		v, ok := vals[name]
		if !ok {
			return 0, false
		}
		if v < 0 {
			return 0, false
		}
		total += v
	}
	return total, true
}
