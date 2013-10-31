// Copyright 2012-2013 Apcera Inc. All rights reserved.

package termtables

import (
	"fmt"
	"strings"
)

type tableAlignment int

// These constants control the alignment which should be used when rendering
// the content of a cell.
const (
	AlignLeft   = tableAlignment(1)
	AlignCenter = tableAlignment(2)
	AlignRight  = tableAlignment(3)
)

// TableStyle controls styling information for a Table as a whole.
//
// For the Border rules, only X, Y and I are needed, and all have defaults.
// The others will all default to the same as BorderI.
type TableStyle struct {
	SkipBorder        bool
	BorderX           string
	BorderY           string
	BorderI           string
	BorderTop         string
	BorderBottom      string
	BorderRight       string
	BorderLeft        string
	BorderTopLeft     string
	BorderTopRight    string
	BorderBottomLeft  string
	BorderBottomRight string
	PaddingLeft       int
	PaddingRight      int
	Width             int
	Alignment         tableAlignment
}

// A CellStyle controls all style applicable to one Cell.
type CellStyle struct {
	// Alignment indicates the alignment to be used in rendering the content
	Alignment tableAlignment

	// ColSpan indicates how many columns this Cell is expected to consume.
	ColSpan int
}

// DefaultStyle is a TableStyle which can be used to get some simple
// default styling for a table, using ASCII characters for drawing borders.
var DefaultStyle = &TableStyle{
	SkipBorder: false,
	BorderX:    "-", BorderY: "|", BorderI: "+",
	PaddingLeft: 1, PaddingRight: 1,
	Width:     80,
	Alignment: AlignLeft,

	// FIXME: the use of a Width here may interact poorly with a changing
	// MaxColumns value; we don't set MaxColumns here because the evaluation
	// order of a var and an init value adds undesired subtlety.
}

type renderStyle struct {
	cellWidths map[int]int
	columns    int

	// used for markdown rendering
	replaceContent func(string) string

	TableStyle
}

// setUtfBoxStyle changes the border characters to be suitable for use when
// the output stream can render UTF-8 characters.
func (s *TableStyle) setUtfBoxStyle() {
	s.BorderX = "─"
	s.BorderY = "│"
	s.BorderI = "┼"
	s.BorderTop = "┬"
	s.BorderBottom = "┴"
	s.BorderLeft = "├"
	s.BorderRight = "┤"
	s.BorderTopLeft = "╭"
	s.BorderTopRight = "╮"
	s.BorderBottomLeft = "╰"
	s.BorderBottomRight = "╯"
}

// setAsciiBoxStyle changes the border characters back to their defaults
func (s *TableStyle) setAsciiBoxStyle() {
	s.BorderX = "-"
	s.BorderY = "|"
	s.BorderI = "+"
	s.BorderTop, s.BorderBottom, s.BorderLeft, s.BorderRight = "", "", "", ""
	s.BorderTopLeft, s.BorderTopRight, s.BorderBottomLeft, s.BorderBottomRight = "", "", "", ""
	s.fillStyleRules()
}

// fillStyleRules populates members of the TableStyle box-drawing specification
// with BorderI as the default.
func (s *TableStyle) fillStyleRules() {
	if s.BorderTop == "" {
		s.BorderTop = s.BorderI
	}
	if s.BorderBottom == "" {
		s.BorderBottom = s.BorderI
	}
	if s.BorderLeft == "" {
		s.BorderLeft = s.BorderI
	}
	if s.BorderRight == "" {
		s.BorderRight = s.BorderI
	}
	if s.BorderTopLeft == "" {
		s.BorderTopLeft = s.BorderI
	}
	if s.BorderTopRight == "" {
		s.BorderTopRight = s.BorderI
	}
	if s.BorderBottomLeft == "" {
		s.BorderBottomLeft = s.BorderI
	}
	if s.BorderBottomRight == "" {
		s.BorderBottomRight = s.BorderI
	}
}

func createRenderStyle(table *Table) *renderStyle {
	style := &renderStyle{TableStyle: *table.Style, cellWidths: map[int]int{}}
	style.TableStyle.fillStyleRules()

	if table.outputMode == outputMarkdown {
		style.buildReplaceContent(table.Style.BorderY)
	}

	// FIXME: handle actually defined width condition

	// loop over the rows and cells to calculate widths
	for _, element := range table.elements {
		// skip separators
		if _, ok := element.(*Separator); ok {
			continue
		}

		// iterate over cells
		if row, ok := element.(*Row); ok {
			totalWidth := 0
			for i, cell := range row.cells {
				totalWidth = totalWidth + cell.Width()
				// FIXME: need to support sizing with colspan handling
				if cell.colSpan > 1 {
					continue
				}
				if style.cellWidths[i] < cell.Width() {
					style.cellWidths[i] = cell.Width()
				}
				if i == len(row.cells)-1 {
					if table.minWidth > totalWidth {
						if style.cellWidths[i] <= cell.Width() {
							// The minus 3 is to avoid odd numbers of padding on right over left-hand side.
							style.cellWidths[i] = cell.Width() + table.minWidth - totalWidth - 3
						}
					}
				}
			}
		}
	}
	style.columns = len(style.cellWidths)

	// calculate actual width
	width := 1 // start at 1 for left border
	for _, v := range style.cellWidths {
		width += v + style.PaddingLeft + style.PaddingRight + 1 // for border
	}
	// right border is covered in loop
	style.Width = width

	return style
}

// CellWidth returns the width of the cell at the supplied index, where the
// width is the number of tty character-cells required to draw the glyphs.
func (s *renderStyle) CellWidth(i int) int {
	return s.cellWidths[i]
}

// buildReplaceContent creates a function closure, with minimal bound lexical
// state, which replaces content
func (s *renderStyle) buildReplaceContent(bad string) {
	replacement := fmt.Sprintf("&#x%02x;", bad)
	s.replaceContent = func(old string) string {
		return strings.Replace(old, bad, replacement, -1)
	}
}
