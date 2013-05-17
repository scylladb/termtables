// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

type tableAlignment int

const (
	AlignLeft   = tableAlignment(1)
	AlignCenter = tableAlignment(2)
	AlignRight  = tableAlignment(3)
)

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

type CellStyle struct {
	Alignment tableAlignment
	ColSpan   int
}

var DefaultStyle = &TableStyle{
	SkipBorder: false,
	BorderX:    "-", BorderY: "|", BorderI: "+",
	PaddingLeft: 1, PaddingRight: 1,
	Width:     80,
	Alignment: AlignLeft,
}

type renderStyle struct {
	cellWidths map[int]int
	columns    int
	TableStyle
}

func (s *TableStyle) setUtfBoxStyle() {
	s.BorderX = "─"
	s.BorderY = "│"
	s.BorderI = "┼"
	s.BorderTop = "┬"
	s.BorderBottom = "┴"
	s.BorderLeft = "├"
	s.BorderRight = "┤"
	s.BorderTopLeft = "┌"
	s.BorderTopRight = "┐"
	s.BorderBottomLeft = "└"
	s.BorderBottomRight = "┘"
}

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

	// FIXME: handle actually defined width condition

	// loop over the rows and cells to calculate widths
	for _, element := range table.elements {
		// skip separators
		if _, ok := element.(*Separator); ok {
			continue
		}

		// iterate over cells
		if row, ok := element.(*Row); ok {
			for i, cell := range row.cells {
				// FIXME: need to support sizing with colspan handling
				if cell.colSpan > 1 {
					continue
				}
				if style.cellWidths[i] < cell.Width() {
					style.cellWidths[i] = cell.Width()
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

func (s *renderStyle) CellWidth(i int) int {
	return s.cellWidths[i]
}
