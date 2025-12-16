package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/kevin-hanselman/slide-puzzle-solver/slide_puzzle"
)

func main() {
	rows := flag.Int("rows", 0, "number of rows in the puzzle")
	cols := flag.Int("cols", 0, "number of columns in the puzzle")
	empty := flag.Int("empty", 0, "value representing the empty tile")
	flag.Parse()

	// Validate flags
	if *rows <= 0 {
		fmt.Fprintf(os.Stderr, "Error: -rows must be positive\n")
		os.Exit(1)
	}
	if *cols <= 0 {
		fmt.Fprintf(os.Stderr, "Error: -cols must be positive\n")
		os.Exit(1)
	}

	// Parse remaining arguments as puzzle values
	args := flag.Args()
	expectedArgs := (*rows) * (*cols)
	if len(args) != expectedArgs {
		fmt.Fprintf(os.Stderr, "Error: expected %d values for %dx%d puzzle, got %d\n", expectedArgs, *rows, *cols, len(args))
		os.Exit(1)
	}

	// Convert string arguments to integers
	values := make([]int, len(args))
	for i, arg := range args {
		val, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid value '%s': %v\n", arg, err)
			os.Exit(1)
		}
		values[i] = val
	}

	// Build grid in row-major order
	grid := make([][]int, *rows)
	idx := 0
	for r := range *rows {
		grid[r] = make([]int, *cols)
		for c := range *cols {
			grid[r][c] = values[idx]
			idx++
		}
	}

	// Create puzzle
	puzzle, err := slide_puzzle.NewPuzzle(grid, *empty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Solve puzzle
	moves, err := puzzle.Solve()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print solution
	if len(moves) == 0 {
		fmt.Println("Puzzle is already solved!")
	} else {
		fmt.Printf("Solution in %d moves:\n", len(moves))
		for i, move := range moves {
			fmt.Printf("%d. %s\n", i+1, move)
		}
	}
}
