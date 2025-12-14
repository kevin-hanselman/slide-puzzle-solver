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

func TestMakeMove(t *testing.T) {
	t.Run("move North - swap with cell above", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		got, err := puzzle.makeMove(North)
		if err != nil {
			t.Fatalf("makeMove(North) error: %v", err)
		}

		wantGrid := [][]int{
			{1, 0, 3},
			{4, 2, 5},
			{6, 7, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		if diff := cmp.Diff(*want, got, cmp.AllowUnexported(Puzzle{}, cell{}, coord{})); diff != "" {
			t.Errorf("makeMove(North) mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("move South - swap with cell below", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		got, err := puzzle.makeMove(South)
		if err != nil {
			t.Fatalf("makeMove(South) error: %v", err)
		}

		wantGrid := [][]int{
			{1, 2, 3},
			{4, 7, 5},
			{6, 0, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		if diff := cmp.Diff(*want, got, cmp.AllowUnexported(Puzzle{}, cell{}, coord{})); diff != "" {
			t.Errorf("makeMove(South) mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("move East - swap with cell to the right", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		got, err := puzzle.makeMove(East)
		if err != nil {
			t.Fatalf("makeMove(East) error: %v", err)
		}

		wantGrid := [][]int{
			{1, 2, 3},
			{4, 5, 0},
			{6, 7, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		if diff := cmp.Diff(*want, got, cmp.AllowUnexported(Puzzle{}, cell{}, coord{})); diff != "" {
			t.Errorf("makeMove(East) mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("move West - swap with cell to the left", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		got, err := puzzle.makeMove(West)
		if err != nil {
			t.Fatalf("makeMove(West) error: %v", err)
		}

		wantGrid := [][]int{
			{1, 2, 3},
			{0, 4, 5},
			{6, 7, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		if diff := cmp.Diff(*want, got, cmp.AllowUnexported(Puzzle{}, cell{}, coord{})); diff != "" {
			t.Errorf("makeMove(West) mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("original puzzle is not modified", func(t *testing.T) {
		originalGrid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(originalGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		// Deep copy the expected grid before the move
		wantGrid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(North)
		if err != nil {
			t.Fatalf("makeMove(North) error: %v", err)
		}

		// Original puzzle should remain unchanged
		if diff := cmp.Diff(*want, *puzzle, cmp.AllowUnexported(Puzzle{}, cell{}, coord{})); diff != "" {
			t.Errorf("original puzzle was modified (-want +got):\n%s", diff)
		}
	})

	t.Run("sequence of moves", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		// Move North, then East
		result, err := puzzle.makeMove(North)
		if err != nil {
			t.Fatalf("makeMove(North) error: %v", err)
		}
		got, err := result.makeMove(East)
		if err != nil {
			t.Fatalf("makeMove(East) error: %v", err)
		}

		wantGrid := [][]int{
			{1, 3, 0},
			{4, 2, 5},
			{6, 7, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		if diff := cmp.Diff(*want, got, cmp.AllowUnexported(Puzzle{}, cell{}, coord{})); diff != "" {
			t.Errorf("after North then East mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("invalid move North from top edge", func(t *testing.T) {
		grid := [][]int{
			{1, 0, 3},
			{4, 2, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(North)
		if err == nil {
			t.Fatal("makeMove(North) from top edge error = nil, want InvalidMoveError")
		}

		var invalidErr *InvalidMoveError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("makeMove(North) from top edge error type = %T, want *InvalidMoveError", err)
		}
	})

	t.Run("invalid move South from bottom edge", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 0, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(South)
		if err == nil {
			t.Fatal("makeMove(South) from bottom edge error = nil, want InvalidMoveError")
		}

		var invalidErr *InvalidMoveError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("makeMove(South) from bottom edge error type = %T, want *InvalidMoveError", err)
		}
	})

	t.Run("invalid move East from right edge", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 5, 0},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(East)
		if err == nil {
			t.Fatal("makeMove(East) from right edge error = nil, want InvalidMoveError")
		}

		var invalidErr *InvalidMoveError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("makeMove(East) from right edge error type = %T, want *InvalidMoveError", err)
		}
	})

	t.Run("invalid move West from left edge", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{0, 4, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(West)
		if err == nil {
			t.Fatal("makeMove(West) from left edge error = nil, want InvalidMoveError")
		}

		var invalidErr *InvalidMoveError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("makeMove(West) from left edge error type = %T, want *InvalidMoveError", err)
		}
	})
}
