package main

import (
	"fmt"
	tl "github.com/JoelOtter/termloop"
	"math"
	"math/rand"
	"strconv"
	"time"
)

////////////////////////
// Maze generation stuff
////////////////////////

type Point struct {
	x int
	y int
	p *Point
}

func (p *Point) Opposite() *Point {
	if p.x != p.p.x {
		return &Point{x: p.x + (p.x - p.p.x), y: p.y, p: p}
	}
	if p.y != p.p.y {
		return &Point{x: p.x, y: p.y + (p.y - p.p.y), p: p}
	}
	return nil
}

func adjacents(point *Point, maze [][]rune) []*Point {
	res := make([]*Point, 0)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if (i == 0 && j == 0) || (i != 0 && j != 0) {
				continue
			}
			if !isInMaze(point.x+i, point.y+j, len(maze), len(maze[0])) {
				continue
			}
			if maze[point.x+i][point.y+j] == '*' {
				res = append(res, &Point{point.x + i, point.y + j, point})
			}
		}
	}
	return res
}

func isInMaze(x, y int, w, h int) bool {
	return x >= 0 && x < w &&
		y >= 0 && y < h
}

// Generates a maze using Prim's Algorithm
// https://en.wikipedia.org/wiki/Maze_generation_algorithm#Randomized_Prim.27s_algorithm
func generateMaze(w, h int) [][]rune {
	maze := make([][]rune, w)
	for row := range maze {
		maze[row] = make([]rune, h)
		for ch := range maze[row] {
			maze[row][ch] = '*'
		}
	}
	rand.NewSource(time.Now().UnixNano())
	point := &Point{x: rand.Intn(w), y: rand.Intn(h)}
	maze[point.x][point.y] = 'S'
	var last *Point
	walls := adjacents(point, maze)
	for len(walls) > 0 {
		rand.NewSource(time.Now().UnixNano())
		wall := walls[rand.Intn(len(walls))]
		for i, w := range walls {
			if w.x == wall.x && w.y == wall.y {
				walls = append(walls[:i], walls[i+1:]...)
				break
			}
		}
		opp := wall.Opposite()
		if isInMaze(opp.x, opp.y, w, h) && maze[opp.x][opp.y] == '*' {
			maze[wall.x][wall.y] = '.'
			maze[opp.x][opp.y] = '.'
			walls = append(walls, adjacents(opp, maze)...)
			last = opp
		}
	}
	maze[last.x][last.y] = 'L'
	bordered := make([][]rune, len(maze)+2)
	for r := range bordered {
		bordered[r] = make([]rune, len(maze[0])+2)
		for c := range bordered[r] {
			if r == 0 || r == len(maze)+1 || c == 0 || c == len(maze[0])+1 {
				bordered[r][c] = '*'
			} else {
				bordered[r][c] = maze[r-1][c-1]
			}
		}
	}
	return bordered
}

/////////////////
// Termloop stuff
/////////////////

type Block struct {
	*tl.Rectangle
	px        int // Previous x
	py        int // Previous y
	move      bool
	g         *tl.Game
	w         int // Width of maze
	h         int // Height of maze
	score     int
	scoretext *tl.Text
	status    int
}

func NewBlock(x, y int, color tl.Attr, g *tl.Game, w, h, score int, scoretext *tl.Text) *Block {
	b := &Block{
		g:         g,
		w:         w,
		h:         h,
		score:     score,
		scoretext: scoretext,
		status:    0,
	}
	b.Rectangle = tl.NewRectangle(x, y, 1, 1, color)
	return b
}

func (b *Block) Draw(s *tl.Screen) {
	if l, ok := b.g.Screen().Level().(*tl.BaseLevel); ok {
		// Set the level offset so the player is always in the
		// center of the screen. This simulates moving the camera.
		sw, sh := s.Size()
		x, y := b.Position()
		l.SetOffset(sw/2-x, sh/2-y)
		b.SetColor(tl.ColorGreen)
	}
	b.Rectangle.Draw(s)
}

func (b *Block) Tick(ev tl.Event) {
	// Enable arrow key movement
	if ev.Type == tl.EventKey {
		b.px, b.py = b.Position()
		switch ev.Key {
		case tl.KeyArrowRight:
			b.SetPosition(b.px+1, b.py)
			b.status--
		case tl.KeyArrowLeft:
			b.SetPosition(b.px-1, b.py)
			b.status--
		case tl.KeyArrowUp:
			b.SetPosition(b.px, b.py-1)
			b.status--
		case tl.KeyArrowDown:
			b.SetPosition(b.px, b.py+1)
			b.status--
		}
	}
}

func (b *Block) Collide(c tl.Physical) {
	if r, ok := c.(*tl.Rectangle); ok {
		switch r.Color() {
		case tl.ColorWhite:
			// Collision with walls
			b.SetPosition(b.px, b.py)
		case tl.ColorBlue:
			// Collision with end - new level!
			b.w *= 2
			b.h *= 2
			b.score *= 2
			buildLevel(b.g, b.w, b.h, b.score)
		case tl.ColorYellow:
			scoreRoot := int(math.Log2(float64(b.score))) - 3
			rand.NewSource(time.Now().UnixNano())
			b.status = 6 + scoreRoot
			for i := 0; i < 5; i++ {
				b.status += rand.Intn(4)
			}
			for i := 0; i < scoreRoot && scoreRoot > 0; i++ {
				switch rand.Intn(32) {
				case 0:
					b.status += rand.Intn(2) + 1
					break
				case 1:
					b.status += rand.Intn(2) + 1
					break
				case 2:
					b.status += rand.Intn(4) + 1
					break
				case 3:
					b.status += rand.Intn(4) + 1
					break
				case 4:
					b.status += rand.Intn(4) + 1
					break
				case 5:
					b.status += rand.Intn(4) + 1
					break
				case 6:
					b.status += rand.Intn(8) + 1
					break
				default:
					b.status++
					break
				}
			}
		case tl.ColorRed:
			if b.status < 1 {
				gameOver()
			}
		}
	}
}

func buildLevel(g *tl.Game, w, h, score int) {
	maze := generateMaze(w, h)
	l := tl.NewBaseLevel(tl.Cell{})
	random := 1
	g.Screen().SetLevel(l)
	g.Log("width: %d / height: %d / Difficulty level: %d", w, h, score)
	g.Screen().AddEntity(tl.NewText(0, 0, "Blue: goal🏆  / Red: dead☠  / Yellow: power up💪 ",
		tl.ColorWhite, tl.ColorBlack))
	scoretext := tl.NewText(0, 2, "Difficulty level: "+strconv.Itoa(score), tl.ColorWhite, tl.ColorBlack)
	g.Screen().AddEntity(scoretext)
	for i, row := range maze {
		for j, path := range row {
			if path == '*' {
				l.AddEntity(tl.NewRectangle(i, j, 1, 1, tl.ColorWhite))
			} else if path == 'S' {
				col := tl.RgbTo256Color(0xff, 0, 0)
				l.AddEntity(NewBlock(i, j, col, g, w, h, score, scoretext))
			} else if path == 'L' {
				l.AddEntity(tl.NewRectangle(i, j, 2, 2, tl.ColorBlue))
			} else if score > 2 {
				switch random {
				case 1:
					l.AddEntity(tl.NewRectangle(i, j, 1, 1, tl.ColorYellow))
				case 2:
					l.AddEntity(tl.NewRectangle(i, j, 1, 1, tl.ColorYellow))
				case 3:
					l.AddEntity(tl.NewRectangle(i, j, 1, 1, tl.ColorRed))
				}
				rand.NewSource(time.Now().UnixNano())
				random = rand.Intn(60)
			}
		}
	}
}

func gameOver() {
	g := tl.NewGame()
	g.Screen().AddEntity(tl.NewText(24, 10, "Game Over", tl.ColorWhite, tl.ColorRed))
	g.Screen().AddEntity(tl.NewText(24, 12, "Press Ctrl+C", tl.ColorWhite, tl.ColorRed))
	g.Start()
}

func main() {
	fmt.Println("終了するには Ctrl+C")
	g := tl.NewGame()
	g.Screen().SetFps(60)
	buildLevel(g, 3, 1, 1)
	g.SetDebugOn(true)
	g.Start()
}
