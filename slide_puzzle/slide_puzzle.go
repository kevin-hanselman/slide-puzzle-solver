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

type InvalidMoveError struct {
	msg string
}

func (e InvalidMoveError) Error() string {
	return e.msg
}

type coord struct {
	row, col int
}

type Move int

const (
	North Move = iota
	South
	East
	West
)

var moveStrings = map[Move]string{
	North: "North",
	South: "South",
	East:  "East",
	West:  "West",
}

func (m Move) String() string {
	return moveStrings[m]
}

func NewPuzzle(grid [][]int, emptyCellValue int) (*Puzzle, error) {
	if len(grid) == 0 {
		return nil, &InvalidPuzzleError{"puzzle must have at least one row"}
	}

	emptyCount := 0
	emptyCell := cell{value: emptyCellValue}
	rowLength := len(grid[0])
	numCells := len(grid) * rowLength
	seen := make([]bool, numCells)

	for row := range grid {
		if len(grid[row]) != rowLength {
			return nil, &InvalidPuzzleError{
				fmt.Sprintf(
					"all rows must have the same length; row 0 has %d columns, row %d has %d columns",
					rowLength,
					row,
					len(grid[row]),
				),
			}
		}
		for col := range grid[row] {
			val := grid[row][col]

			// Check if value is in valid range
			if val < 0 || val >= numCells {
				return nil, &InvalidPuzzleError{fmt.Sprintf("grid values must be in range [0, %d); got %d", numCells, val)}
			}

			// Check for duplicates
			if seen[val] {
				return nil, &InvalidPuzzleError{fmt.Sprintf("duplicate value %d found in grid", val)}
			}
			seen[val] = true

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

func (p Puzzle) getMoves() map[Move]bool {
	moves := make(map[Move]bool)
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

func (p Puzzle) isSolved() bool {
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

// makeMove moves the Puzzle's emptyCell in the given direction and returns the
// updated Puzzle.
//
// Note that the receiver is not a pointer, so the original puzzle is not
// modified.
func (p Puzzle) makeMove(m Move) (Puzzle, error) {
	// Validate that the move is possible
	validMoves := p.getMoves()
	if !validMoves[m] {
		return Puzzle{}, &InvalidMoveError{fmt.Sprintf("cannot move %s from current position", m)}
	}

	// Create a deep copy of the grid
	newGrid := make([][]int, len(p.grid))
	for i := range p.grid {
		newGrid[i] = make([]int, len(p.grid[i]))
		copy(newGrid[i], p.grid[i])
	}

	// Determine the target cell based on move direction
	var targetRow, targetCol int
	switch m {
	case North:
		targetRow = p.emptyCell.coord.row - 1
		targetCol = p.emptyCell.coord.col
	case South:
		targetRow = p.emptyCell.coord.row + 1
		targetCol = p.emptyCell.coord.col
	case East:
		targetRow = p.emptyCell.coord.row
		targetCol = p.emptyCell.coord.col + 1
	case West:
		targetRow = p.emptyCell.coord.row
		targetCol = p.emptyCell.coord.col - 1
	}

	// Swap the empty cell with the target cell
	newGrid[p.emptyCell.coord.row][p.emptyCell.coord.col] = newGrid[targetRow][targetCol]
	newGrid[targetRow][targetCol] = p.emptyCell.value

	// Create new puzzle with updated grid and empty cell position
	return Puzzle{
		grid: newGrid,
		emptyCell: cell{
			value: p.emptyCell.value,
			coord: coord{row: targetRow, col: targetCol},
		},
	}, nil
}
