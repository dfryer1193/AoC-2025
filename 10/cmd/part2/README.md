# Advent of Code 2025 - Day 10, Part 2

This program solves the second part of the Day 10 puzzle from Advent of Code 2025. It calculates the minimum number of button presses required to configure the joltage levels of a series of machines.

## Running the Solution

There are two ways to run this solution: using the native Go solvers or using a more optimized version accelerated with an external ILP (Integer Linear Programming) solver.

### Standard Version (Native Go Solver)

This version uses built-in solvers written in pure Go. It requires no external dependencies beyond the Go standard library.

To run the standard version, use the following command:

```bash
go run ./ <input_file>
```

Replace `<input_file>` with the path to your puzzle input file (e.g., `input.txt` or `test.txt`).

Example:
```bash
go run ./ test.txt
```

### Accelerated Version (with `golp` and `lpsolve`)

This version uses the `github.com/draffensperger/golp` package, which is a Go wrapper for the `lpsolve` library. This can be significantly faster for complex inputs but requires `lpsolve` to be installed on your system.

#### 1. Install `lpsolve`

You must install the `lpsolve` library and its development headers. The installation command depends on your operating system:

-   **Arch Linux:**
    ```bash
    sudo pacman -S lpsolve
    ```
-   **Debian/Ubuntu:**
    ```bash
    sudo apt-get install lp-solve
    ```
-   **macOS (using Homebrew):**
    ```bash
    brew install lp_solve
    ```

#### 2. Run with `golp` build tag

To compile and run the accelerated version, use the `-tags golp` build flag:

```bash
go run -tags golp ./ <input_file>
```

Example:
```bash
go run -tags golp ./ input.txt
```

#### Troubleshooting `lpsolve` Installation

If the Go compiler cannot find the `lpsolve` header files (e.g., you see an error like `fatal error: lp_lib.h: No such file or directory`), you may need to provide the paths to the compiler and linker manually using `CGO_CFLAGS` and `CGO_LDFLAGS`.

For example, if `lp_lib.h` is in `/usr/include/lpsolve` and the library is in `/usr/lib`, you can use the following command:

```bash
CGO_CFLAGS="-I/usr/include/lpsolve" CGO_LDFLAGS="-L/usr/lib -llpsolve55" go run -tags golp ./ input.txt
```

Adjust the paths according to where `lpsolve` was installed on your system.
