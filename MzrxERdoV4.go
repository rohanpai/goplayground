package main

import (
	&#34;image&#34;
	&#34;image/draw&#34;
)

type Mouse struct {
	Loc	image.Point
	Buttons	int
}

// Context represents the context for a GUI client.
type Context struct {
	// W receives an value when the window changes.
	W	&lt;-chan *Window

	// K receives an value when a key a pressed.
	K	&lt;-chan rune

	// M receives an value when the mouse moves.
	M	&lt;-chan Mouse

	// Dragging receives a value when another window
	// is dragging something into this context.
	Dragging	&lt;-chan DragEvent
}

type DragEventKind int

const (
	_	DragEventKind	= iota
	DragEnter
	DragMove
	DragDrop
	DragLeave
)

// A DragEvent notifies a GUI client of an object being dragged
// onto it.
type DragEvent struct {
	Kind	DragEventKind
	// Loc gives the current position of the dragged object.
	Loc	image.Point
	// Data holds the data associated with the dragged object.
	Data	interface{}
	// On a DragDrop event, Reply will be non-nil
	// and must be used to indicate whether the dropped
	// object has been accepted.
	Reply	chan&lt;- bool
}

// Window represents a GUI window.
type Window struct {
	// TODO
}

func (w *Window) Image() draw.Image {
	// TODO
	return nil
}

func (w *Window) Flush() {
	// TODO
}

func (w *Window) Drag(data interface{}, icon image.Image, from image.Point) bool {
	// TODO read from mouse channel and distribute dragging events
	// to appropriate clients.
	return false
}

// ----------------------------------------------
// Some sample client code.
// The client has a window with multiple slots, each
// of which can contain an item. The user can drag
// items between slots in different windows.

type SomeClient struct {
	slots []*Slot
}

type Content struct {
	s string
}

// image returns an image for the being-dragged item.
func (c *Content) image() image.Image {
	// TODO
	return nil
}

type Slot struct {
	// r holds the rectangle in the window covered by the slot.
	r	image.Rectangle

	// content holds the current contents of the slot.
	content	*Content

	// willing signifies that this slot is showing readiness to accept a
	// dragged object.
	willing	bool
}

func NewContext() *Context {
	return &amp;Context{
	// TODO
	}
}

func (c *SomeClient) Run() {
	ctxt := NewContext()
	win := &lt;-ctxt.W
	for {
		select {
		case win = &lt;-ctxt.W:
			// window replaced.
		case m := &lt;-ctxt.M:
			if slot := c.SlotAtPoint(m.Loc); slot != nil &amp;&amp; slot.content != nil {
				if win.Drag(slot.content, slot.content.image(), m.Loc) {
					c.UpdateSlot(slot, nil)
				}
			}
		case e := &lt;-ctxt.Dragging:
			// Some other window is dragging an object into this one.
			c.dragging(e, ctxt.Dragging)
		}
	}
}

func (c *SomeClient) dragging(e DragEvent, dragc &lt;-chan DragEvent) {
	var hover *Slot	// slot the drag is currently hovering over.
	for {
		content, _ := e.Data.(*Content)
		slot := c.SlotAtPoint(e.Loc)
		switch e.Kind {
		case DragLeave:
			return
		case DragMove:
			// ignore unknown drags
			if content == nil || slot == hover {
				break
			}
			if hover != nil {
				c.ShowWilling(hover, false)
			}
			if slot != nil {
				c.ShowWilling(slot, true)
			}
			hover = slot

		case DragDrop:
			if hover != nil {
				c.ShowWilling(hover, false)
			}
			if content == nil || slot == nil {
				e.Reply &lt;- false
				return
			}
			// Accept the drop and complete the drag-and-drop action.
			c.UpdateSlot(slot, content)
			e.Reply &lt;- true
			return
		}
		e = &lt;-dragc
	}
}

func (c *SomeClient) SlotAtPoint(p image.Point) *Slot {
	// TODO find slot under given point.
	return c.slots[0]
}

func (c *SomeClient) ShowWilling(slot *Slot, willing bool) {
	// TODO update displayed image to reflect new status.
	slot.willing = willing
}

func (c *SomeClient) UpdateSlot(slot *Slot, content *Content) {
	// TODO update displayed image to reflect new content
}

func main() {
}
