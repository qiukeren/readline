package core

import (
	"github.com/reeflective/readline/internal/strutil"
	"github.com/xo/inputrc"
)

// Cursor is the cursor position in the current line buffer.
// Contains methods to set, move, describe and check itself.
type Cursor struct {
	pos  int
	mark int
	line *Line
}

// NewCursor is a required constructor for the line cursor,
// because some default numeric values must be negative.
func NewCursor(line *Line) *Cursor {
	return &Cursor{
		pos:  0,
		mark: -1,
		line: line,
	}
}

// Set sets the position of the cursor to an absolute value.
// If either negative or greater than the length of the line,
// the cursor will be set to either 0, or the length of the line.
func (c *Cursor) Set(pos int) {
	defer c.CheckAppend()

	switch {
	case pos < 0:
		c.pos = 0
	case pos > c.line.Len():
		c.pos = c.line.Len()
	default:
		c.pos = pos
	}
}

// Pos returns the current cursor position.
// This function cannot return an invalid cursor position: it cannot be negative, nor it
// can be greater than the length of the line (note it still can be out of line by 1).
func (c *Cursor) Pos() int {
	c.CheckAppend()
	return c.pos
}

// Inc increments the cursor position by 1, if it's not at the end of the line.
func (c *Cursor) Inc() {
	if c.pos < c.line.Len() {
		c.pos++
	}
}

// Dec decrements the cursor position by 1, if it's not at the beginning of the line.
func (c *Cursor) Dec() {
	if c.pos > 0 {
		c.pos--
	}
}

// Move moves the cursor position by a relative value. If the end result is negative,
// the cursor is set to 0. If longer than the line, the cursor is set to length of line.
func (c *Cursor) Move(offset int) {
	defer c.CheckAppend()
	c.pos += offset
}

// ToFirstNonSpace moves the cursor either backward or forward to
// the first character in the line that is not a space, a tab or
// a newline. If the current is not one, the cursor doesn't move.
// If the cursor is at the end of the line, the move is performed
// backward, regardless of the forward parameter value.
func (c *Cursor) ToFirstNonSpace(forward bool) {
	if c.line.Len() == 0 {
		return
	}

	defer c.CheckAppend()

	if c.pos == c.line.Len() {
		forward = false
		c.pos--
	}

	// At line bounds
	if forward && c.pos >= c.line.Len()-1 {
		return
	} else if !forward && c.pos == 0 {
		return
	}

	for {
		if !c.onSpace() {
			return
		}

		if forward {
			c.pos++
		} else {
			c.pos--
		}

		if c.pos <= 0 && c.pos >= c.line.Len() {
			return
		}

	}
}

// BeginningOfLine moves the cursor to the beginning of the current line,
// (marked by a newline) or if no newline found, to the beginning of the buffer.
func (c *Cursor) BeginningOfLine() {
	defer c.CheckCommand()

	newlinePos := c.line.Find(inputrc.Newline, c.pos, false)
	if newlinePos != -1 {
		c.pos = newlinePos + 1
	} else {
		c.pos = 0
	}
}

// EndOfLine moves the cursor to the end of the current line, (marked by
// a newline) or if no newline found, to the position of the last character.
func (c *Cursor) EndOfLine() {
	defer c.CheckCommand()

	if c.OnEmptyLine() {
		return
	}

	newlinePos := c.line.Find(inputrc.Newline, c.pos, true)

	if newlinePos != -1 {
		c.pos = newlinePos - 1
	} else {
		c.pos = c.line.Len() - 1
	}
}

// EndOfLineAppend moves the cursor to the very end of the line,
// that is, equal to len(Line), as in when appending in insert mode.
func (c *Cursor) EndOfLineAppend() {
	defer c.CheckAppend()

	newlinePos := c.line.Find(inputrc.Newline, c.pos, true)

	if newlinePos != -1 {
		c.pos = newlinePos
	} else {
		c.pos = c.line.Len()
	}
}

// SetMark sets the current cursor position as the mark.
func (c *Cursor) SetMark() {
	c.CheckAppend()
	c.mark = c.pos
}

// Mark returns the current mark value of the cursor, or -1 if not set.
func (c *Cursor) Mark() int {
	return c.mark
}

// ResetMark resets the insertion point mark.
func (c *Cursor) ResetMark() {
	c.mark = -1
}

// Line returns the index of the current line on which the cursor is.
// A line is defined as a sequence of runes between one or two newline
// characters, between end and/or beginning of buffer, or a mix of both.
func (c *Cursor) Line() int {
	c.CheckAppend()

	newlines := c.line.newlines()

	// Either match between two newlines
	for i, newline := range newlines {
		if newline[0] < c.pos {
			continue
		}

		return i
	}

	// Or return the number of lines
	return len(newlines)
}

// LineMove moves the cursor by n lines either up (if the value is negative),
// or down (if positive). If greater than the length of possible lines above/below,
// the cursor will be set to either the first, or the last line of the buffer.
func (c *Cursor) LineMove(lines int) {
	c.CheckAppend()
	defer c.CheckAppend()

	newlines := c.line.newlines()
	if len(newlines) == 1 || lines == 0 {
		return
	}

	if lines < 0 {
		for i := 0; i < -1*lines; i++ {
			c.moveLineUp()
			c.CheckCommand()
		}
	} else {
		for i := 0; i < lines; i++ {
			c.moveLineDown()
			c.CheckCommand()
		}
	}
}

// OnEmptyLine returns true if the rune under the current cursor position is a newline
// and that the preceding rune in the line is also a newline, or returns false.
func (c *Cursor) OnEmptyLine() bool {
	if c.line.Len() == 0 {
		return true
	}

	if c.pos == 0 {
		return (*c.line)[c.pos] == inputrc.Newline
	} else if c.pos == c.line.Len() {
		return (*c.line)[c.pos-1] == inputrc.Newline
	}

	return (*c.line)[c.pos] == inputrc.Newline
}

// AtBeginningOfLine returns true if the cursor is either at the beginning
// of the line buffer, or on the first character after the previous newline.
func (c *Cursor) AtBeginningOfLine() bool {
	if c.pos == 0 {
		return false
	}

	newlines := c.line.newlines()

	for line := 0; line < len(newlines); line++ {
		epos := newlines[line][0]
		if epos == c.pos-1 {
			return true
		}
	}

	return false
}

// CheckAppend verifies that the current cursor position is neither negative,
// nor greater than the length of the input line. If either is true, the
// cursor will set its value as either 0, or the length of the line.
func (c *Cursor) CheckAppend() {
	// Position
	if c.pos < 0 {
		c.pos = 0
	}

	if c.pos > c.line.Len() {
		c.pos = c.line.Len()
	}

	// Mark, invalid position deactivates it.
	if c.mark < -1 {
		c.mark = -1
	}

	if c.mark > c.line.Len() {
		c.mark = -1
	}
}

// CheckCommand is like CheckAppend, but ensures the cursor position is never greater
// than the length of the line minus 1, since in command mode, the cursor is on a char.
func (c *Cursor) CheckCommand() {
	c.CheckAppend()

	if c.pos == c.line.Len() {
		c.pos--
	}

	// The cursor can also not be on a newline sign,
	// as it will induce the line rendering into an error.
	if c.line.Len() > 0 && (*c.line)[c.pos] == '\n' && !c.OnEmptyLine() {
		c.Inc()
	}
}

// Coordinates returns the number of real terminal lines above the cursor position
// (y value), and the number of columns since the beginning of the current line (x value).
// Params:
// @indent -    Used to align all lines (except the first) together on a single column.
func (c *Cursor) Coordinates(indent int) (x, y int) {
	c.CheckAppend()

	newlines := c.line.newlines()
	bpos := 0
	usedY := 0

	for pos, newline := range newlines {
		switch {
		case newline[0] < c.pos:
			// Until we didn't reach the cursor line,
			// simply care about the line count.
			line := (*c.line)[bpos:newline[0]]
			bpos = newline[0] + 1
			_, y := strutil.LineSpan(line, pos, indent)
			usedY += y

		default:
			// On the cursor line, use both line and column count.
			line := (*c.line)[bpos:c.pos]
			usedX, y := strutil.LineSpan(line, pos, indent)
			usedY += y

			return usedX, usedY
		}
	}

	return
}

func (c *Cursor) moveLineDown() {
	var cpos, begin int

	newlines := c.line.newlines()

	for line := 0; line < len(newlines); line++ {
		end := newlines[line][0]
		if line < c.Line() {
			begin = end
			continue
		}

		// If we are on the current line,
		// go at the end of it
		if line == c.Line() {
			cpos = c.pos - begin
			begin = end
			continue
		}

		// And either go at the end of the line
		// or to the previous cursor X coordinate.
		if end-begin > cpos {
			c.pos = begin + cpos
		} else {
			c.pos = end - 1
		}

		break
	}
}

func (c *Cursor) moveLineUp() {
	var cpos, begin int

	newlines := c.line.newlines()

	for line := len(newlines) - 1; line >= 0; line-- {
		end := newlines[line][0] - 1

		if line > c.Line() {
			continue
		}

		// Get the beginning of the previous line.
		if line > 0 {
			begin = newlines[line-1][0] + 1
		} else {
			begin = 0
			end--
		}

		// If we are on the current line,
		// go at the beginning of the previous one.
		if line == c.Line() {
			cpos = c.pos - begin
			continue
		}

		// And either go at the end of the line
		// or to the previous cursor X coordinate.
		if end-begin > cpos {
			c.pos = begin + cpos
		} else {
			c.pos = end
		}

		break
	}
}

func (c *Cursor) onSpace() bool {
	switch (*c.line)[c.pos] {
	case inputrc.Space, inputrc.Newline, inputrc.Tab:
		return true
	default:
		return false
	}
}
