package slide_puzzle

import "fmt"

type Puzzle struct {
	grid      [][]int
	emptyTile tile
}

type tile struct {
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

func NewPuzzle(grid [][]int, emptyTileValue int) (*Puzzle, error) {
	if len(grid) == 0 {
		return nil, &InvalidPuzzleError{"puzzle must have at least one row"}
	}

	emptyCount := 0
	emptyTile := tile{value: emptyTileValue}
	rowLength := len(grid[0])
	numTiles := len(grid) * rowLength
	seen := make([]bool, numTiles)

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
			if val < 0 || val >= numTiles {
				return nil, &InvalidPuzzleError{fmt.Sprintf("grid values must be in range [0, %d); got %d", numTiles, val)}
			}

			// Check for duplicates
			if seen[val] {
				return nil, &InvalidPuzzleError{fmt.Sprintf("duplicate value %d found in grid", val)}
			}
			seen[val] = true

			if grid[row][col] == emptyTileValue {
				emptyCount++
				emptyTile.coord = coord{row: row, col: col}
			}
		}
	}

	if emptyCount != 1 {
		return nil, &InvalidPuzzleError{fmt.Sprintf("puzzle must contain exactly one empty tile; got %d", emptyCount)}
	}

	return &Puzzle{grid: grid, emptyTile: emptyTile}, nil
}

func (p Puzzle) getMoves() map[Move]bool {
	moves := make(map[Move]bool)
	// North: move tile from south up
	if p.emptyTile.coord.row < (len(p.grid) - 1) {
		moves[North] = true
	}
	// South: move tile from north down
	if p.emptyTile.coord.row > 0 {
		moves[South] = true
	}
	// East: move tile from west right
	if p.emptyTile.coord.col > 0 {
		moves[East] = true
	}
	// West: move tile from east left
	if p.emptyTile.coord.col < (len(p.grid[0]) - 1) {
		moves[West] = true
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

// makeMove moves a tile in the given direction into the empty space and returns
// the updated Puzzle. For example, Move North moves the tile south of the empty
// space up into the empty space.
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

	// Determine the target tile based on move direction
	// Move direction refers to the tile moving, not the empty tile.
	var targetRow, targetCol int
	switch m {
	case North:
		// Move tile from south up
		targetRow = p.emptyTile.coord.row + 1
		targetCol = p.emptyTile.coord.col
	case South:
		// Move tile from north down
		targetRow = p.emptyTile.coord.row - 1
		targetCol = p.emptyTile.coord.col
	case East:
		// Move tile from west right
		targetRow = p.emptyTile.coord.row
		targetCol = p.emptyTile.coord.col - 1
	case West:
		// Move tile from east left
		targetRow = p.emptyTile.coord.row
		targetCol = p.emptyTile.coord.col + 1
	}

	// Swap the empty tile with the target tile.
	newGrid[p.emptyTile.coord.row][p.emptyTile.coord.col] = newGrid[targetRow][targetCol]
	newGrid[targetRow][targetCol] = p.emptyTile.value

	return Puzzle{
		grid: newGrid,
		emptyTile: tile{
			value: p.emptyTile.value,
			coord: coord{row: targetRow, col: targetCol},
		},
	}, nil
}
