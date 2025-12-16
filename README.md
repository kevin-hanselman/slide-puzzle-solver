# Slide Puzzle Solver

Solves sliding tile puzzles using breadth-first search.

Written as another experiment with vibe coding.

## Usage

```bash
go run main.go -rows <n> -cols <m> -empty <value> <tile1> <tile2> ... <tileN>
```

- Tiles are specified in row-major order (left-to-right, top-to-bottom).
- Move directions describe which tile moves into the empty space (e.g., "North" moves the tile below the empty space upward)
- Uses BFS to find the shortest solution
- Goal state: tiles arranged sequentially from `0` to `n-1`
