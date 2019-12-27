package ebitest

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type (
	// Dir represents a direction.
	Dir        int
	mouseState int
	touchState int
	// Control represents a direction.
	Control int
)

const (
	// DirUp ...
	DirUp Dir = iota
	// DirRight ...
	DirRight
	// DirDown ...
	DirDown
	// DirLeft ...
	DirLeft
)

const (
	mouseStateNone mouseState = iota
	mouseStatePressing
	mouseStateSettled
)

const (
	touchStateNone touchState = iota
	touchStatePressing
	touchStateSettled
	touchStateInvalid
)

const (
	// CtrlA ...
	CtrlA Control = iota
	// CtrlZ ...
	CtrlZ
)

// String returns a string representing the direction.
func (d Dir) String() string {
	switch d {
	case DirUp:
		return "Up"
	case DirRight:
		return "Right"
	case DirDown:
		return "Down"
	case DirLeft:
		return "Left"
	}
	panic("not reach")
}

// Vector returns a [-1, 1] value for each axis.
func (d Dir) Vector() (x, y int) {
	switch d {
	case DirUp:
		return 0, -1
	case DirRight:
		return 1, 0
	case DirDown:
		return 0, 1
	case DirLeft:
		return -1, 0
	}
	panic("not reach")
}

// String returns a string representing the direction.
func (c Control) String() string {
	switch c {
	case CtrlA:
		return "A"
	case CtrlZ:
		return "Z"
	}
	panic("not reach")
}

// Input represents the current key states.
type Input struct {
	mouseState    mouseState
	mouseInitPosX int
	mouseInitPosY int
	mouseDir      Dir

	touchState    touchState
	touchID       int
	touchInitPosX int
	touchInitPosY int
	touchLastPosX int
	touchLastPosY int
	touchDir      Dir
}

// NewInput generates a new Input object.
func NewInput() *Input {
	return &Input{}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func vecToDir(dx, dy int) (Dir, bool) {
	if abs(dx) < 4 && abs(dy) < 4 {
		return 0, false
	}
	if abs(dx) < abs(dy) {
		if dy < 0 {
			return DirUp, true
		}
		return DirDown, true
	}
	if dx < 0 {
		return DirLeft, true
	}
	return DirRight, true
}

// Update updates the current input states.
func (i *Input) Update() {
	switch i.mouseState {
	case mouseStateNone:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			i.mouseInitPosX = x
			i.mouseInitPosY = y
			i.mouseState = mouseStatePressing
		}
	case mouseStatePressing:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			dx := x - i.mouseInitPosX
			dy := y - i.mouseInitPosY
			d, ok := vecToDir(dx, dy)
			if !ok {
				i.mouseState = mouseStateNone
				break
			}
			i.mouseDir = d
			i.mouseState = mouseStateSettled
		}
	case mouseStateSettled:
		i.mouseState = mouseStateNone
	}
	switch i.touchState {
	case touchStateNone:
		ts := ebiten.TouchIDs()
		if len(ts) == 1 {
			i.touchID = ts[0]
			x, y := ebiten.TouchPosition(ts[0])
			i.touchInitPosX = x
			i.touchInitPosY = y
			i.touchLastPosX = x
			i.touchLastPosX = y
			i.touchState = touchStatePressing
		}
	case touchStatePressing:
		ts := ebiten.TouchIDs()
		if len(ts) >= 2 {
			break
		}
		if len(ts) == 1 {
			if ts[0] != i.touchID {
				i.touchState = touchStateInvalid
			} else {
				x, y := ebiten.TouchPosition(ts[0])
				i.touchLastPosX = x
				i.touchLastPosY = y
			}
			break
		}
		if len(ts) == 0 {
			dx := i.touchLastPosX - i.touchInitPosX
			dy := i.touchLastPosY - i.touchInitPosY
			d, ok := vecToDir(dx, dy)
			if !ok {
				i.touchState = touchStateNone
				break
			}
			i.touchDir = d
			i.touchState = touchStateSettled
		}
	case touchStateSettled:
		i.touchState = touchStateNone
	case touchStateInvalid:
		if len(ebiten.TouchIDs()) == 0 {
			i.touchState = touchStateNone
		}
	}
}

// Dir returns a currently pressed direction.
// Dir returns false if no direction key is pressed.
func (i *Input) Dir() (Dir, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		return DirUp, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		return DirLeft, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		return DirRight, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		return DirDown, true
	}
	if i.mouseState == mouseStateSettled {
		return i.mouseDir, true
	}
	if i.touchState == touchStateSettled {
		return i.touchDir, true
	}
	return 0, false
}

// Control returns false if no control key is pressed.
func (i *Input) Control() (Control, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		return CtrlA, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		return CtrlZ, true
	}
	return 0, false
}
