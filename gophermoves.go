package gophermoves

import (
	"bytes"
	"os"
	"strings"

	"golang.org/x/term"
)

// Moves is an interface that defines methods to simulate movements for the character like Gopher
type Moves interface {

	// Start will initiate the play to simulate movements for the character
	Start() (err error)

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

	fd int

	move  chan struct{}
	close chan struct{}

	term *term.State
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

	QuitLowerCase = []byte(`q`)
	QuitUpperCase = []byte(`Q`)

	QuitUsingCtrlC = []byte{3}

	initMsg       = "Move your char: [W,S,A,D] == [UP,DOWN,LEFT,RIGHT]\r\n"
	cleanTerminal = "\033[H\033[2J"
)

// New will create and return a new instance that satisfies the Moves interface for Gopher
func New(m int) Moves {
	e := &implMoves{
		positionX: 0,
		positionY: 0,
		matrix:    m,
		move:      make(chan struct{}, 1),
		close:     make(chan struct{}, 1),
		term:      nil,
	}

	e.fd = int(os.Stdin.Fd())

	return e
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
	_, _ = os.Stdin.WriteString(msg)
}

func (e *implMoves) scan() {
	var b = make([]byte, 1)
	for {
		_, _ = os.Stdin.Read(b)
		switch {
		case bytes.Equal(b, UpLowerCase) || bytes.Equal(b, UpUpperCase):
			e.Up()
			e.move <- struct{}{}
		case bytes.Equal(b, DownLowerCase) || bytes.Equal(b, DownUpperCase):
			e.Down()
			e.move <- struct{}{}
		case bytes.Equal(b, LeftLowerCase) || bytes.Equal(b, LeftUpperCase):
			e.Left()
			e.move <- struct{}{}
		case bytes.Equal(b, RightLowerCase) || bytes.Equal(b, RightUpperCase):
			e.Right()
			e.move <- struct{}{}
		case bytes.Equal(b, ResetLowerCase) || bytes.Equal(b, ResetUpperCase):
			e.Reset()
			e.move <- struct{}{}
		case bytes.Equal(b, QuitLowerCase) || bytes.Equal(b, QuitUpperCase):
			e.close <- struct{}{}
		case bytes.Equal(b, QuitUsingCtrlC):
			e.close <- struct{}{}
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
		msgShow += "\r\n"
	}

	e.print(msgShow)
}

func (e *implMoves) Reset() {
	e.positionX = 0
	e.positionY = 0
}

func (e *implMoves) Start() (err error) {
	e.makeTerm()
	e.show()

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

func (e *implMoves) makeTerm() {
	t, err := term.MakeRaw(e.fd)
	if err != nil {
		panic(err)
	}
	e.term = t
}
