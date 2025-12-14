package slide_puzzle

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewPuzzle(t *testing.T) {
	t.Run("valid puzzle with one empty cell", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}

		got, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error = %v, want nil", err)
		}

		if diff := cmp.Diff(grid, got.grid); diff != "" {
			t.Errorf("NewPuzzle() grid mismatch (-want +got):\n%s", diff)
		}

		if got.emptyCell.coord.row != 1 || got.emptyCell.coord.col != 1 {
			t.Errorf("NewPuzzle() emptyCell.coord = {row: %d, col: %d}, want {row: 1, col: 1}",
				got.emptyCell.coord.row, got.emptyCell.coord.col)
		}

		if got.emptyCell.value != 0 {
			t.Errorf("NewPuzzle() emptyCell.value = %d, want 0", got.emptyCell.value)
		}
	})

	t.Run("empty cell at top-left", func(t *testing.T) {
		grid := [][]int{
			{0, 1, 2},
			{3, 4, 5},
			{6, 7, 8},
		}

		got, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error = %v, want nil", err)
		}

		if got.emptyCell.coord.row != 0 || got.emptyCell.coord.col != 0 {
			t.Errorf("NewPuzzle() emptyCell.coord = {row: %d, col: %d}, want {row: 0, col: 0}",
				got.emptyCell.coord.row, got.emptyCell.coord.col)
		}
	})

	t.Run("empty cell at bottom-right", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 0},
		}

		got, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error = %v, want nil", err)
		}

		if got.emptyCell.coord.row != 2 || got.emptyCell.coord.col != 2 {
			t.Errorf("NewPuzzle() emptyCell.coord = {row: %d, col: %d}, want {row: 2, col: 2}",
				got.emptyCell.coord.row, got.emptyCell.coord.col)
		}
	})

	t.Run("multiple empty cells returns error", func(t *testing.T) {
		grid := [][]int{
			{1, 0, 3},
			{4, 0, 5},
			{6, 7, 8},
		}

		_, err := NewPuzzle(grid, 0)
		if err == nil {
			t.Fatal("NewPuzzle() error = nil, want InvalidPuzzleError")
		}

		var invalidErr *InvalidPuzzleError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("NewPuzzle() error type = %T, want *InvalidPuzzleError", err)
		}
	})

	t.Run("no empty cells returns error", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}

		_, err := NewPuzzle(grid, 0)
		if err == nil {
			t.Fatal("NewPuzzle() error = nil, want InvalidPuzzleError")
		}

		var invalidErr *InvalidPuzzleError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("NewPuzzle() error type = %T, want *InvalidPuzzleError", err)
		}
	})

	t.Run("jagged rows return error", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0},
			{6, 7, 8},
		}

		_, err := NewPuzzle(grid, 0)
		if err == nil {
			t.Fatal("NewPuzzle() error = nil, want InvalidPuzzleError")
		}

		var invalidErr *InvalidPuzzleError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("NewPuzzle() error type = %T, want *InvalidPuzzleError", err)
		}
	})
}

func TestGetMove(t *testing.T) {
	tests := []struct {
		name      string
		grid      [][]int
		wantMoves map[move]bool
	}{
		{
			name: "empty cell in middle - all moves available",
			grid: [][]int{
				{1, 2, 3},
				{4, 0, 5},
				{6, 7, 8},
			},
			wantMoves: map[move]bool{North: true, South: true, East: true, West: true},
		},
		{
			name: "empty cell at top-left corner",
			grid: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			wantMoves: map[move]bool{South: true, East: true},
		},
		{
			name: "empty cell at top-right corner",
			grid: [][]int{
				{1, 2, 0},
				{3, 4, 5},
				{6, 7, 8},
			},
			wantMoves: map[move]bool{South: true, West: true},
		},
		{
			name: "empty cell at bottom-left corner",
			grid: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{0, 7, 8},
			},
			wantMoves: map[move]bool{North: true, East: true},
		},
		{
			name: "empty cell at bottom-right corner",
			grid: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 0},
			},
			wantMoves: map[move]bool{North: true, West: true},
		},
		{
			name: "empty cell on top edge",
			grid: [][]int{
				{1, 0, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			wantMoves: map[move]bool{South: true, East: true, West: true},
		},
		{
			name: "empty cell on left edge",
			grid: [][]int{
				{1, 2, 3},
				{0, 4, 5},
				{6, 7, 8},
			},
			wantMoves: map[move]bool{North: true, South: true, East: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			puzzle, err := NewPuzzle(tt.grid, 0)
			if err != nil {
				t.Fatalf("NewPuzzle() error: %v", err)
			}

			got := puzzle.getMoves()

			if diff := cmp.Diff(tt.wantMoves, got); diff != "" {
				t.Errorf("getMoves() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIsSolved(t *testing.T) {
	tests := []struct {
		name   string
		grid   [][]int
		expect bool
	}{
		{
			name: "solved 3x3 puzzle",
			grid: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			expect: true,
		},
		{
			name: "unsolved 3x3 puzzle - swapped tiles",
			grid: [][]int{
				{0, 1, 2},
				{3, 5, 4},
				{6, 7, 8},
			},
			expect: false,
		},
		{
			name: "solved 2x2 puzzle",
			grid: [][]int{
				{0, 1},
				{2, 3},
			},
			expect: true,
		},
		{
			name: "unsolved 2x2 puzzle",
			grid: [][]int{
				{1, 0},
				{2, 3},
			},
			expect: false,
		},
		{
			name: "solved 4x3 puzzle",
			grid: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
				{9, 10, 11},
			},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			puzzle, err := NewPuzzle(tt.grid, 0)
			if err != nil {
				t.Fatalf("NewPuzzle() error: %v", err)
			}

			got := puzzle.isSolved()

			if got != tt.expect {
				t.Errorf("isSolved() = %v, want %v", got, tt.expect)
			}
		})
	}
}
