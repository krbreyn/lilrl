package game

import (
	"math/rand"
	"slices"
)

// definitely need to make this whole file cleaner after i figure this stuff out and work on
// the rest of the featureset

func makeNewGrid(width, height int) [][]Tile {
	grid := make([][]Tile, height)
	for i := range grid {
		grid[i] = make([]Tile, width)
	}
	for y := range grid {
		for x := range grid[y] {
			grid[y][x] = NewTile(TileWall)
		}
	}
	return grid
}

func GenNewFloor(depth int) ([][]Tile, []Vec2) {
	if depth == 1 {
		grid := makeNewGrid(10, 10)
		for y := range grid {
			for x := range grid[y] {
				if y == 0 || x == 0 || y == 9 || x == 9 {
					continue
				}
				grid[y][x] = NewTile(TileFloor)
			}
		}
		grid[rand.Intn(8)+1][rand.Intn(8)+1] = NewTile(TileStairDown)
		return grid, []Vec2{Vec2{5, 5}}
	}
	return GenDngRogue(100, 50)
}

type Section struct {
	ID               int // 0..8 for a 3x3 grid
	Row, Col         int // position in section grid
	RoomX, RoomY     int // top-left coordinate
	RoomWidth        int
	RoomHeight       int
	CenterX, CenterY int // center coord of room
	Connected        bool
	ConnectedIDs     []int
}

func divideGridIntoSections(grid [][]Tile, n int) []Section {
	totalRows := len(grid)
	if totalRows == 0 {
		return nil
	}
	totalCols := len(grid[0])

	sectionHeight := totalRows / n
	sectionWidth := totalCols / n

	var sections []Section
	sectionID := 0

	minRoomWidth, minRoomHeight := 3, 3

	for i := range n {
		for j := range n {
			secRowStart := i * sectionHeight
			secColStart := j * sectionWidth

			maxRoomWidth := sectionWidth - 2
			maxRoomHeight := sectionHeight - 2

			roomWidth := minRoomWidth + rand.Intn(maxRoomWidth-minRoomWidth+1)
			roomHeight := minRoomHeight + rand.Intn(maxRoomHeight-minRoomHeight+1)

			roomXOffset := 1 + rand.Intn(sectionWidth-roomWidth-1)
			roomYOffset := 1 + rand.Intn(sectionHeight-roomHeight-1)

			roomX := secColStart + roomXOffset
			roomY := secRowStart + roomYOffset

			centerX := roomX + roomWidth/2
			centerY := roomY + roomHeight/2

			sec := Section{
				ID:           sectionID,
				Row:          i,
				Col:          j,
				RoomX:        roomX,
				RoomY:        roomY,
				RoomWidth:    roomWidth,
				RoomHeight:   roomHeight,
				CenterX:      centerX,
				CenterY:      centerY,
				Connected:    false,
				ConnectedIDs: []int{},
			}
			sections = append(sections, sec)
			sectionID++
		}
	}

	return sections
}

// getNeighbors returns the IDs of adjacent sections (up, down, left, right)
// for the section at index id in a grid with n sections per row.
func getNeighbors(sections []Section, id, n int) []int {
	var neighbors []int
	sec := sections[id]
	row, col := sec.Row, sec.Col
	// Up
	if row > 0 {
		neighbors = append(neighbors, (row-1)*n+col)
	}
	// Down
	if row < n-1 {
		neighbors = append(neighbors, (row+1)*n+col)
	}
	// Left
	if col > 0 {
		neighbors = append(neighbors, row*n+col-1)
	}
	// Right
	if col < n-1 {
		neighbors = append(neighbors, row*n+col+1)
	}
	return neighbors
}

func drawRoom(grid [][]Tile, sec Section) {
	for y := sec.RoomY; y < sec.RoomY+sec.RoomHeight; y++ {
		for x := sec.RoomX; x < sec.RoomX+sec.RoomWidth; x++ {
			grid[y][x] = NewTile(TileFloor)
		}
	}
}

func drawLCorridor(grid [][]Tile, x1, y1, x2, y2 int) {
	flip := rand.Intn(2)
	if flip == 0 {
		// L-shape: horizontal then vertical
		if x1 < x2 {
			for x := x1; x <= x2; x++ {
				grid[y1][x] = NewTile(TileFloor)
			}
		} else {
			for x := x2; x <= x1; x++ {
				grid[y1][x] = NewTile(TileFloor)
			}
		}
		if y1 < y2 {
			for y := y1; y <= y2; y++ {
				grid[y][x2] = NewTile(TileFloor)
			}
		} else {
			for y := y2; y <= y1; y++ {
				grid[y][x2] = NewTile(TileFloor)
			}
		}
	} else {
		// L-shape: vertical then horizontal
		if y1 < y2 {
			for y := y1; y <= y2; y++ {
				grid[y][x1] = NewTile(TileFloor)
			}
		} else {
			for y := y2; y <= y1; y++ {
				grid[y][x1] = NewTile(TileFloor)
			}
		}
		if x1 < x2 {
			for x := x1; x <= x2; x++ {
				grid[y2][x] = NewTile(TileFloor)
			}
		} else {
			for x := x2; x <= x1; x++ {
				grid[y2][x] = NewTile(TileFloor)
			}
		}
	}
}

func drawBCorridor(grid [][]Tile, x1, y1, x2, y2 int) {
	path := Bresenham(x1, y1, x2, y2)
	for _, dir := range path {
		grid[y1][x1] = NewTile(TileFloor)
		x1 += dir.X - x1
		y1 += dir.Y - y1
	}
}

func isInRoom(x, y int, sections []Section) bool {
	for _, sec := range sections {
		if x >= sec.RoomX && x < sec.RoomX+sec.RoomWidth &&
			y >= sec.RoomY && y < sec.RoomY+sec.RoomHeight {
			return true
		}
	}
	return false
}

func isAdjacentToRoom(x, y int, sections []Section) bool {
	if isInRoom(x, y-1, sections) {
		return true
	}
	if isInRoom(x, y+1, sections) {
		return true
	}
	if isInRoom(x-1, y, sections) {
		return true
	}
	if isInRoom(x+1, y, sections) {
		return true
	}
	return false
}

// func floodFillCount(grid [][]Tile, x, y int, visited map[[2]int]bool, min int) int {
// 	queue := [][2]int{{x, y}}
// 	count := 0

// 	for len(queue) > 0 {
// 		if count >= min {
// 			break
// 		}
// 		current := queue[0]
// 		queue = queue[1:]

// 		cx, cy := current[0], current[1]
// 		if visited[[2]int{cx, cy}] {
// 			continue
// 		}
// 		visited[[2]int{cx, cy}] = true
// 		count++

// 		directions := [][2]int{
// 			{0, -1},
// 			{0, 1},
// 			{-1, 0},
// 			{1, 0},
// 		}

// 		for _, d := range directions {
// 			nx, ny := cx+d[0], cy+d[1]
// 			if nx >= 0 && nx < len(grid[0]) && ny >= 0 && ny < len(grid) &&
// 				grid[ny][nx].Type == TileFloor && !visited[[2]int{nx, ny}] {
// 				queue = append(queue, [2]int{nx, ny})
// 			}
// 		}
// 	}

// 	return count
// }

func GenCaveRandomWalk(width, height int) [][]Tile {
	grid := makeNewGrid(width, height)

	moveRandom := func(x, y int) (int, int) {
		directions := []struct{ dx, dy int }{
			{0, -1}, // Up
			{0, 1},  // Down
			{-1, 0}, // Left
			{1, 0},  // Right
			//{-1, -1}, // Up-Left
			//{1, -1},  // Up-Right
			//{1, 1},   // Down-Right
			//{-1, 1},  // Down-Let
		}

		for {
			dir := directions[rand.Intn(len(directions))]
			newX, newY := x+dir.dx, y+dir.dy

			if newX > 0 && newX < width-1 && newY > 0 && newY < height-1 {
				return newX, newY
			} else {
				/* reduce the frequency of ugly maps by randomly replacing the random walker when it
				reaches the edge of the map */
				for {
					testx, testy := rand.Intn(width-2)+1, rand.Intn(height-2)+1
					if grid[testy][testx].Type == TileFloor {
						return testx, testy
					}
				}
			}
		}
	}

	// pick starting point
	pos_y := rand.Intn(height-2) + 1
	pos_x := rand.Intn(width-2) + 1
	grid[pos_y][pos_x] = NewTile(TileFloor)

	numTiles := width * height
	numFloor := 0

	modifier := ((rand.Float64()*2 - 1) * 5) / 100 // +/- 5%
	desiredFloorPercent := 0.45 + modifier

	for (float64(numFloor) / float64(numTiles)) < desiredFloorPercent {
		pos_x, pos_y = moveRandom(pos_x, pos_y)
		if grid[pos_y][pos_x].Type == TileWall {
			grid[pos_y][pos_x] = NewTile(TileFloor)
			numFloor++
		}
	}

	return grid
}

// randomly fill the [][]tile with walls/spaces and then smooth
// it out in passes converting every tile where a certain amount of tiles in the
// 3x3 area are walls (twiddle with it to see what works)
func GenCaveAutomata() {

}

func GenMazeBacktrack() {

}

func GenMazePrims() {

}

// TODO integrate random walk into dng generation by having
// a random walker randomly map out the room interior of a section

func GenDngRogue(width, height int) ([][]Tile, []Vec2) {
	grid := makeNewGrid(width, height)
	n := rand.Intn(2) + 5
	sections := divideGridIntoSections(grid, n)
	totalSections := len(sections)

	// connect sections using DFS
	visited := make([]bool, totalSections)
	stack := []int{rand.Intn(totalSections)}
	visited[stack[0]] = true

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		neighbors := getNeighbors(sections, current, n)
		rand.Shuffle(len(neighbors), func(i, j int) {
			neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
		})

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				sections[current].ConnectedIDs = append(sections[current].ConnectedIDs, neighbor)
				sections[neighbor].ConnectedIDs = append(sections[neighbor].ConnectedIDs, current)
				visited[neighbor] = true
				stack = append(stack, neighbor)
			}
		}
	}

	//add extra random connections
	extraConnections := rand.Intn(n+1) + n
	for range extraConnections {
		a := rand.Intn(totalSections)
		neighbors := getNeighbors(sections, a, n)
		if len(neighbors) == 0 {
			continue
		}

		b := neighbors[rand.Intn(len(neighbors))]
		// avoid duplicate connections
		alreadyConnected := slices.Contains(sections[a].ConnectedIDs, b)
		if !alreadyConnected {
			sections[a].ConnectedIDs = append(sections[a].ConnectedIDs, b)
			sections[b].ConnectedIDs = append(sections[b].ConnectedIDs, a)
		}
	}

	for _, sec := range sections {
		drawRoom(grid, sec)
	}

	for _, sec := range sections {
		for _, connID := range sec.ConnectedIDs {
			other := sections[connID]
			flip := rand.Intn(3)
			if flip == 0 {
				drawLCorridor(grid, sec.CenterX, sec.CenterY, other.CenterX, other.CenterY)
			} else {
				drawBCorridor(grid, sec.CenterX, sec.CenterY, other.CenterX, other.CenterY)
			}
		}
	}

	// scan for doors
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			if grid[y][x].Type == TileFloor {
				wallCount := 0
				if grid[y-1][x].Type == TileWall {
					wallCount++
				}
				if grid[y+1][x].Type == TileWall {
					wallCount++
				}
				if grid[y][x-1].Type == TileWall {
					wallCount++
				}
				if grid[y][x+1].Type == TileWall {
					wallCount++
				}
				if wallCount >= 2 && !isInRoom(x, y, sections) && isAdjacentToRoom(x, y, sections) {
					grid[y][x] = NewTile(TileDoor)
				}
			}
		}
	}

	// place stairs
	adj := make([][]int, totalSections)
	for i, sec := range sections {
		adj[i] = sec.ConnectedIDs
	}
	upstairs := placeStairs(adj, sections, grid, 3)

	return grid, upstairs
}

// Computer distances from a set of starting sections.
func bfsDistances(adj [][]int, startSet []int) []int {
	n := len(adj)
	distances := make([]int, n)
	for i := range distances {
		distances[i] = -1 // unvisited
	}

	queue := make([]int, 0, n)
	for _, s := range startSet {
		distances[s] = 0
		queue = append(queue, s)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, neighbor := range adj[current] {
			if distances[neighbor] == -1 {
				distances[neighbor] = distances[current] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return distances
}

// Place k pairs of stairs
func placeStairs(adj [][]int, sections []Section, grid [][]Tile, k int) []Vec2 {
	upStairs := []int{rand.Intn(len(sections))} // first up stair
	downStairs := []int{}

	for i := range k {
		// place down stair farthest from the current up stair
		distances := bfsDistances(adj, []int{upStairs[i]})
		maxDist := -1
		downStair := -1
		for j, dist := range distances {
			if dist > maxDist && !slices.Contains(upStairs, j) && !slices.Contains(downStairs, j) {
				maxDist = dist
				downStair = j
			}
		}
		if downStair == -1 {
			break // no more valid sections
		}
		downStairs = append(downStairs, downStair)

		if i < k-1 { // place next up stair if not the last pair
			allStairs := append(upStairs, downStairs...)
			distances = bfsDistances(adj, allStairs)
			maxDist = -1
			nextUpStair := -1
			for j, dist := range distances {
				if dist > maxDist && !slices.Contains(upStairs, j) && !slices.Contains(downStairs, j) {
					maxDist = dist
					nextUpStair = j
				}
			}
			if nextUpStair == -1 {
				break
			}
			upStairs = append(upStairs, nextUpStair)
		}

	}

	var ret []Vec2
	// place stairs
	for _, up := range upStairs {
		y := sections[up].CenterY
		x := sections[up].CenterX
		grid[y][x] = NewTile(TileStairUp)
		ret = append(ret, Vec2{x, y})
	}
	for _, down := range downStairs {
		grid[sections[down].CenterY][sections[down].CenterX] = NewTile(TileStairDown)
	}

	return ret
}

// should be just like the above but with rooms randomly placed instead of
// within divided grid sections
func GenDngAngbang() {

}

// // BFS helper
// func findFarthestSection(adj [][]int, start int) int {
// 	queue := []int{start}
// 	visited := make([]bool, len(adj))
// 	visited[start] = true
// 	farthest := start

// 	for len(queue) > 0 {
// 		current := queue[0]
// 		queue = queue[1:]
// 		farthest = current

// 		for _, neighbor := range adj[current] {
// 			if !visited[neighbor] {
// 				visited[neighbor] = true
// 				queue = append(queue, neighbor)
// 			}
// 		}
// 	}

// 	return farthest
// }
