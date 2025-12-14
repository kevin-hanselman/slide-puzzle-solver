package slide_puzzle

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertPuzzlesEqual(t *testing.T, a, b *Puzzle) {
	if diff := cmp.Diff(a, b, cmp.AllowUnexported(Puzzle{}, tile{}, coord{})); diff != "" {
		t.Errorf("Puzzles are different (-first +second):\n%s", diff)
	}
}

func TestNewPuzzle(t *testing.T) {
	t.Run("valid puzzle with one empty tile", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 8},
		}

		got, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error = %v, want nil", err)
		}

		want := &Puzzle{
			grid: [][]int{
				{1, 2, 3},
				{4, 0, 5},
				{6, 7, 8},
			},
			emptyTile: tile{
				value: 0,
				coord: coord{row: 1, col: 1},
			},
		}

		assertPuzzlesEqual(t, want, got)
	})

	t.Run("empty tile at top-left", func(t *testing.T) {
		grid := [][]int{
			{0, 1, 2},
			{3, 4, 5},
			{6, 7, 8},
		}

		got, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error = %v, want nil", err)
		}

		want := &Puzzle{
			grid: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			emptyTile: tile{
				value: 0,
				coord: coord{row: 0, col: 0},
			},
		}

		assertPuzzlesEqual(t, want, got)
	})

	t.Run("empty tile at bottom-right", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 0},
		}

		got, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error = %v, want nil", err)
		}

		want := &Puzzle{
			grid: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 0},
			},
			emptyTile: tile{
				value: 0,
				coord: coord{row: 2, col: 2},
			},
		}

		assertPuzzlesEqual(t, want, got)
	})

	t.Run("multiple empty tiles returns error", func(t *testing.T) {
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

	t.Run("no empty tiles returns error", func(t *testing.T) {
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

	t.Run("duplicate values return error", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 5, 8},
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

	t.Run("missing values return error", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 9},
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

	t.Run("values outside valid range return error", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, 10},
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

	t.Run("negative values return error", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 0, 5},
			{6, 7, -1},
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
		wantMoves map[Move]bool
	}{
		{
			name: "empty tile in middle - all moves available",
			grid: [][]int{
				{1, 2, 3},
				{4, 0, 5},
				{6, 7, 8},
			},
			wantMoves: map[Move]bool{North: true, South: true, East: true, West: true},
		},
		{
			name: "empty tile at top-left corner",
			grid: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			wantMoves: map[Move]bool{North: true, West: true},
		},
		{
			name: "empty tile at top-right corner",
			grid: [][]int{
				{1, 2, 0},
				{3, 4, 5},
				{6, 7, 8},
			},
			wantMoves: map[Move]bool{North: true, East: true},
		},
		{
			name: "empty tile at bottom-left corner",
			grid: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{0, 7, 8},
			},
			wantMoves: map[Move]bool{South: true, West: true},
		},
		{
			name: "empty tile at bottom-right corner",
			grid: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 0},
			},
			wantMoves: map[Move]bool{South: true, East: true},
		},
		{
			name: "empty tile on top edge",
			grid: [][]int{
				{1, 0, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			wantMoves: map[Move]bool{North: true, East: true, West: true},
		},
		{
			name: "empty tile on left edge",
			grid: [][]int{
				{1, 2, 3},
				{0, 4, 5},
				{6, 7, 8},
			},
			wantMoves: map[Move]bool{North: true, South: true, West: true},
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
	t.Run("move North - tile from south moves up", func(t *testing.T) {
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
			{1, 2, 3},
			{4, 7, 5},
			{6, 0, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		assertPuzzlesEqual(t, want, &got)
	})

	t.Run("move South - tile from north moves down", func(t *testing.T) {
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
			{1, 0, 3},
			{4, 2, 5},
			{6, 7, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		assertPuzzlesEqual(t, want, &got)
	})

	t.Run("move East - tile from west moves right", func(t *testing.T) {
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
			{0, 4, 5},
			{6, 7, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		assertPuzzlesEqual(t, want, &got)
	})

	t.Run("move West - tile from east moves left", func(t *testing.T) {
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
			{4, 5, 0},
			{6, 7, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		assertPuzzlesEqual(t, want, &got)
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
		assertPuzzlesEqual(t, want, puzzle)
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

		// Move North (7 moves up), then East (6 moves right)
		result, err := puzzle.makeMove(North)
		if err != nil {
			t.Fatalf("makeMove(North) error: %v", err)
		}
		got, err := result.makeMove(East)
		if err != nil {
			t.Fatalf("makeMove(East) error: %v", err)
		}

		wantGrid := [][]int{
			{1, 2, 3},
			{4, 7, 5},
			{0, 6, 8},
		}
		want, err := NewPuzzle(wantGrid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		assertPuzzlesEqual(t, want, &got)
	})

	t.Run("invalid move North from bottom edge", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 0, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(North)
		if err == nil {
			t.Fatal("makeMove(North) from bottom edge error = nil, want InvalidMoveError")
		}

		var invalidErr *InvalidMoveError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("makeMove(North) from bottom edge error type = %T, want *InvalidMoveError", err)
		}
	})

	t.Run("invalid move South from top edge", func(t *testing.T) {
		grid := [][]int{
			{1, 0, 3},
			{4, 2, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(South)
		if err == nil {
			t.Fatal("makeMove(South) from top edge error = nil, want InvalidMoveError")
		}

		var invalidErr *InvalidMoveError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("makeMove(South) from top edge error type = %T, want *InvalidMoveError", err)
		}
	})

	t.Run("invalid move East from left edge", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{0, 4, 5},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(East)
		if err == nil {
			t.Fatal("makeMove(East) from left edge error = nil, want InvalidMoveError")
		}

		var invalidErr *InvalidMoveError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("makeMove(East) from left edge error type = %T, want *InvalidMoveError", err)
		}
	})

	t.Run("invalid move West from right edge", func(t *testing.T) {
		grid := [][]int{
			{1, 2, 3},
			{4, 5, 0},
			{6, 7, 8},
		}
		puzzle, err := NewPuzzle(grid, 0)
		if err != nil {
			t.Fatalf("NewPuzzle() error: %v", err)
		}

		_, err = puzzle.makeMove(West)
		if err == nil {
			t.Fatal("makeMove(West) from right edge error = nil, want InvalidMoveError")
		}

		var invalidErr *InvalidMoveError
		if !errors.As(err, &invalidErr) {
			t.Fatalf("makeMove(West) from right edge error type = %T, want *InvalidMoveError", err)
		}
	})
}
