// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

import "strings"

type Row struct {
	cells []*Cell
}

func CreateRow(items []interface{}) *Row {
	row := &Row{ cells: []*Cell{} }
	for i, item := range items {
		if c, ok := item.(*Cell); ok {
			c.column = i
			row.cells = append(row.cells, c)
		} else {
			row.cells = append(row.cells, createCell(i, item, nil))
		}
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
