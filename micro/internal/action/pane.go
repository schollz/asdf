package action

import (
	"github.com/schollz/asdf/micro/internal/display"
)

// A Pane is a general interface for a window in the editor.
type Pane interface {
	Handler
	display.Window
	ID() uint64
	SetID(i uint64)
	Name() string
	Close()
	SetTab(t *Tab)
	Tab() *Tab
}
