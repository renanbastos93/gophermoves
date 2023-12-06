package gophermoves

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

// Moves is an interface that defines methods to simulate movements for the character like Gopher
type Moves interface {

	// Start will initiate the play to simulate movements for the character
	Start()

	// Reset will redefine the default states of the X and Y positions
	Reset()

	// Up moves the character upward
	Up()

	// Down moves the character down
	Down()

	// Left will turn the character to the left
	Left()

	// Right will turn the character to the right
	Right()
}

// implMoves implementation of a structure that satisfies the Moves interface
type implMoves struct {
	positionX int
	positionY int

	matrix int

	in  *bufio.Scanner
	out *bufio.Writer

	move  chan struct{}
	close chan struct{}
}

const (
	square  = "▩"
	hashtag = "#"
)

var (
	UpLowerCase    = []byte(`w`)
	DownLowerCase  = []byte(`s`)
	LeftLowerCase  = []byte(`a`)
	RightLowerCase = []byte(`d`)
	ResetLowerCase = []byte(`r`)

	UpUpperCase    = []byte(`W`)
	DownUpperCase  = []byte(`S`)
	LeftUpperCase  = []byte(`A`)
	RightUpperCase = []byte(`D`)
	ResetUpperCase = []byte(`R`)

	initMsg       = "Move your char: [W,S,A,D] == [UP,DOWN,LEFT,RIGHT]\n"
	cleanTerminal = "\033[H\033[2J"
)

// New will create and return a new instance that satisfies the Moves interface for Gopher
func New(m int) Moves {
	return &implMoves{
		positionX: 0,
		positionY: 0,
		matrix:    m,
		in:        bufio.NewScanner(os.Stdin),
		out:       bufio.NewWriter(os.Stdout),
		move:      make(chan struct{}, 1),
		close:     make(chan struct{}, 1),
	}
}

func (e *implMoves) Up() {
	e.positionY -= 1
}

func (e *implMoves) Down() {
	e.positionY += 1
}

func (e *implMoves) Right() {
	e.positionX += 1
}

func (e *implMoves) Left() {
	e.positionX -= 1
}

func (e implMoves) print(msg string) {
	_, _ = e.out.WriteString(msg)
	_ = e.out.Flush()
}

func (e *implMoves) scan() {
	e.show()
	s := e.in
	for e.in.Scan() {
		switch {
		case bytes.Equal(s.Bytes(), UpLowerCase) || bytes.Equal(s.Bytes(), UpUpperCase):
			e.Up()
			e.move <- struct{}{}
		case bytes.Equal(s.Bytes(), DownLowerCase) || bytes.Equal(s.Bytes(), DownUpperCase):
			e.Down()
			e.move <- struct{}{}
		case bytes.Equal(s.Bytes(), LeftLowerCase) || bytes.Equal(s.Bytes(), LeftUpperCase):
			e.Left()
			e.move <- struct{}{}
		case bytes.Equal(s.Bytes(), RightLowerCase) || bytes.Equal(s.Bytes(), RightUpperCase):
			e.Right()
			e.move <- struct{}{}
		case bytes.Equal(s.Bytes(), ResetLowerCase) || bytes.Equal(s.Bytes(), ResetUpperCase):
			e.Reset()
			e.move <- struct{}{}
		}
	}
}

func (e *implMoves) limitPositions(wh [][]string) [][]string {
	switch {
	case e.positionX > e.matrix-1:
		e.positionX = e.matrix - 1
	case e.positionX < 0:
		e.positionX = 0
	case e.positionY > e.matrix-1:
		e.positionY = e.matrix - 1
	case e.positionY < 0:
		e.positionY = 0
	}

	wh[e.positionY][e.positionX] = square
	return wh
}

func (e *implMoves) setPosition() [][]string {
	var wh = make([][]string, e.matrix)
	for i := 0; i < e.matrix; i++ {
		wh[i] = make([]string, e.matrix)
		for j := 0; j < e.matrix; j++ {
			wh[i][j] = hashtag
		}
	}

	return wh
}

func (e *implMoves) show() {
	msgShow := cleanTerminal
	msgShow += initMsg

	wh := e.setPosition()
	wh = e.limitPositions(wh)
	for i := 0; i < len(wh); i++ {
		msgShow += strings.Join(wh[i], "  ")
		msgShow += "\n"
	}

	e.print(msgShow)
}

func (e *implMoves) Reset() {
	e.positionX = 0
	e.positionY = 0
}

func (e *implMoves) Start() {
	go e.scan()
	for {
		select {
		case <-e.move:
			e.show()
		case <-e.close:
			return
		}
	}
}
