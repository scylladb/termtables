// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

import "strings"

type Row struct {
	cells []*Cell
}

func CreateRow(items []interface{}) *Row {
	row := &Row{cells: []*Cell{}}
	for _, item := range items {
		row.AddCell(item)
	}
	return row
}

func (r *Row) AddCell(item interface{}) {
	if c, ok := item.(*Cell); ok {
		c.column = len(r.cells)
		r.cells = append(r.cells, c)
	} else {
		r.cells = append(r.cells, createCell(len(r.cells), item, nil))
	}
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
