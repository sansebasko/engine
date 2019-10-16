package window

import (
	"time"
)

// Millis
const DefaultDoubleClickDuration = 300

// MouseState keeps track of the state of pressed mouse buttons.
type MouseState struct {
	win                 IWindow
	lastButton          MouseButton
	DoubleClickDuration time.Duration
	states              map[MouseButton]*mouseButtonState
}

type mouseButtonState struct {
	clickCount int
	lastClick  time.Time
}

func (s *mouseButtonState) doubleClicked() bool {
	return s.clickCount == 2 || s.clickCount == -2
}

// NewMouseState returns a new MouseState object.
func NewMouseState(win IWindow) *MouseState {

	ms := new(MouseState)
	ms.win = win
	ms.DoubleClickDuration = DefaultDoubleClickDuration * time.Millisecond
	ms.states = map[MouseButton]*mouseButtonState{
		MouseButtonLeft:   {clickCount: 0, lastClick: time.Now()},
		MouseButtonRight:  {clickCount: 0, lastClick: time.Now()},
		MouseButtonMiddle: {clickCount: 0, lastClick: time.Now()},
	}

	// Subscribe to window mouse events
	ms.win.SubscribeID(OnMouseUp, &ms, ms.onMouseUp)
	ms.win.SubscribeID(OnMouseDown, &ms, ms.onMouseDown)

	return ms
}

// Dispose unsubscribes from the window events.
func (ms *MouseState) Dispose() {

	ms.win.UnsubscribeID(OnMouseUp, &ms)
	ms.win.UnsubscribeID(OnMouseDown, &ms)
}

// Pressed returns whether the specified mouse button is currently pressed.
func (ms *MouseState) Pressed(b MouseButton) bool {

	return ms.states[b].clickCount > 0
}

// Pressed returns whether the left mouse button is currently pressed.
func (ms *MouseState) LeftPressed() bool {

	return ms.states[MouseButtonLeft].clickCount > 0
}

// Pressed returns whether the right mouse button is currently pressed.
func (ms *MouseState) RightPressed() bool {

	return ms.states[MouseButtonRight].clickCount > 0
}

// Pressed returns whether the middle mouse button is currently pressed.
func (ms *MouseState) MiddlePressed() bool {

	return ms.states[MouseButtonMiddle].clickCount > 0
}

// Pressed returns whether the user left double clicked.
func (ms *MouseState) LeftDoubleClicked() bool {

	return ms.lastButton == MouseButtonLeft && ms.states[MouseButtonLeft].doubleClicked()
}

// Pressed returns whether the user right double clicked.
func (ms *MouseState) RightDoubleClicked() bool {

	return ms.lastButton == MouseButtonRight && ms.states[MouseButtonRight].doubleClicked()
}

// Pressed returns whether the user middle double clicked.
func (ms *MouseState) MiddleDoubleClicked() bool {

	return ms.lastButton == MouseButtonMiddle && ms.states[MouseButtonMiddle].doubleClicked()
}

// onMouse receives mouse events and updates the internal map of states.
func (ms *MouseState) onMouseUp(evname string, ev interface{}) {

	mev := ev.(*MouseEvent)
	if ms.states[mev.Button].clickCount > 0 {
		ms.states[mev.Button].clickCount *= -1
	}
}

// onMouse receives mouse events and updates the internal map of states.
func (ms *MouseState) onMouseDown(evname string, ev interface{}) {

	mev := ev.(*MouseEvent)
	ms.lastButton = mev.Button

	now := time.Now()

	if ms.states[mev.Button].clickCount == 0 {
		ms.states[mev.Button].clickCount = 1
		ms.states[mev.Button].lastClick = now
		return
	}

	if ms.states[mev.Button].clickCount == -1 {
		if ms.states[mev.Button].lastClick.Add(ms.DoubleClickDuration).Before(now) {
			ms.states[mev.Button].clickCount = 1
			ms.states[mev.Button].lastClick = now
			return
		}

		ms.states[mev.Button].clickCount = 2
		return
	}

	ms.states[mev.Button].clickCount = 1
	ms.states[mev.Button].lastClick = now
}
