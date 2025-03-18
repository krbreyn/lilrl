package game

import "math/rand"

func GenCaveRandomWalk(width, height int) [][]Tile {
	grid := make([][]Tile, height)
	for i := range grid {
		grid[i] = make([]Tile, width)
	}
	for y := range grid {
		for x := range grid[y] {
			grid[y][x] = NewTile(TileWall)
		}
	}

	moveRandom := func(x, y int) (int, int) {
		directions := []struct{ dx, dy int }{
			{0, -1},  // Up
			{0, 1},   // Down
			{-1, 0},  // Left
			{1, 0},   // Right
			{-1, -1}, // Up-Left
			{1, -1},  // Up-Right
			{1, 1},   // Down-Right
			{-1, 1},  // Down-Let
		}

		for {
			dir := directions[rand.Intn(len(directions))]
			newX, newY := x+dir.dx, y+dir.dy

			if newX > 0 && newX < width-1 && newY > 0 && newY < height-1 {
				return newX, newY
			}
		}
	}

	// pick starting point
	pos_y := rand.Intn(height-2) + 1
	pos_x := rand.Intn(width-2) + 1
	grid[pos_y][pos_x] = NewTile(TileFloor)

	numTiles := width * height
	numFloor := 0

	desiredFloorPercent := 0.65

	for (float64(numFloor) / float64(numTiles)) < desiredFloorPercent {
		pos_x, pos_y = moveRandom(pos_x, pos_y)
		if grid[pos_y][pos_x].Type == TileWall {
			grid[pos_y][pos_x] = NewTile(TileFloor)
			numFloor++
		}
	}

	return grid
}
