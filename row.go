// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

import "strings"

type Row struct {
	cells []*Cell
}

func CreateRow(items []interface{}) *Row {
	row := &Row{ cells: []*Cell{} }
	for i, item := range items {
		row.cells = append(row.cells, CreateCell(i, item))
	}
	return row
}

func (r *Row) Render(style *renderStyle) string {
	// pre-render and shove into an array... helps with cleanly adding borders
	renderedCells := []string{}
	for _, c := range r.cells {
		renderedCells = append(renderedCells, c.Render(style))
	}

	// format final output
	return style.BorderY + strings.Join(renderedCells, style.BorderY) + style.BorderY
}
