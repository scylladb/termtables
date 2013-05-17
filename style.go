// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

type tableAlignment int

const (
	AlignLeft   = tableAlignment(1)
	AlignCenter = tableAlignment(2)
	AlignRight  = tableAlignment(3)
)

type TableStyle struct {
	SkipBorder   bool
	BorderX      string
	BorderY      string
	BorderI      string
	PaddingLeft  int
	PaddingRight int
	Width        int
	Alignment    tableAlignment
}

type CellStyle struct {
	Alignment tableAlignment
	ColSpan   int
}

var DefaultStyle = &TableStyle{SkipBorder: false, BorderX: "-", BorderY: "|",
	BorderI: "+", PaddingLeft: 1, PaddingRight: 1, Width: 80, Alignment: AlignLeft}

type renderStyle struct {
	cellWidths map[int]int
	columns    int
	TableStyle
}

func createRenderStyle(table *Table) *renderStyle {
	style := &renderStyle{TableStyle: *table.Style, cellWidths: map[int]int{}}

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
