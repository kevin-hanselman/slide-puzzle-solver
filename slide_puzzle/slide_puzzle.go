package slide_puzzle

import "fmt"

type Puzzle struct {
	grid      [][]int
	emptyCell cell
}

type cell struct {
	value int
	coord coord
}

type InvalidPuzzleError struct {
	msg string
}

func (e InvalidPuzzleError) Error() string {
	return e.msg
}

type coord struct {
	row, col int
}

type move int

const (
	North move = iota
	South
	East
	West
)

var moveStrings = map[move]string{
	North: "North",
	South: "South",
	East:  "East",
	West:  "West",
}

func (m move) String() string {
	return moveStrings[m]
}

func NewPuzzle(grid [][]int, emptyCellValue int) (*Puzzle, error) {
	if len(grid) == 0 {
		return nil, &InvalidPuzzleError{"puzzle must have at least one row"}
	}

	emptyCount := 0
	emptyCell := cell{value: emptyCellValue}
	rowLength := len(grid[0])

	for row := range grid {
		if len(grid[row]) != rowLength {
			return nil, &InvalidPuzzleError{fmt.Sprintf("all rows must have the same length; row 0 has %d columns, row %d has %d columns", rowLength, row, len(grid[row]))}
		}
		for col := range grid[row] {
			if grid[row][col] == emptyCellValue {
				emptyCount++
				emptyCell.coord = coord{row: row, col: col}
			}
		}
	}

	if emptyCount != 1 {
		return nil, &InvalidPuzzleError{fmt.Sprintf("puzzle must contain exactly one empty cell; got %d", emptyCount)}
	}

	return &Puzzle{grid: grid, emptyCell: emptyCell}, nil
}

func (p *Puzzle) getMoves() map[move]bool {
	moves := make(map[move]bool)
	if p.emptyCell.coord.row > 0 {
		moves[North] = true
	}
	if p.emptyCell.coord.row < (len(p.grid) - 1) {
		moves[South] = true
	}
	if p.emptyCell.coord.col > 0 {
		moves[West] = true
	}
	if p.emptyCell.coord.col < (len(p.grid[0]) - 1) {
		moves[East] = true
	}
	return moves
}

func (p *Puzzle) isSolved() bool {
	want := 0
	for row := range p.grid {
		for col := range p.grid[row] {
			if p.grid[row][col] != want {
				return false
			}
			want++
		}
	}
	return true
}
