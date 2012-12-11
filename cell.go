// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

import (
	"math"
	"strconv"
	"strings"
)

type Cell struct {
	column int
	formattedValue string
	alignment *tableAlignment
	colSpan int
}

func CreateCell(v interface{}, style *CellStyle) *Cell {
	return createCell(0, v, style)
}

func createCell(column int, v interface{}, style *CellStyle) *Cell {
	cell := &Cell{ column: column, formattedValue: renderValue(v), colSpan: 1 }
	if style != nil {
		cell.alignment = &style.Alignment
		if style.ColSpan != 0 {
			cell.colSpan = style.ColSpan
		}
	}
	return cell
}

func (c *Cell) Width() int {
	return len(c.formattedValue)
}

func (c *Cell) Render(style *renderStyle) (buffer string) {
	// if no alignment is set, import the table's default
	if c.alignment == nil {
		c.alignment = &style.Alignment
	}

	// left padding
	buffer += strings.Repeat(" ", style.PaddingLeft)

	// append the main value and handle alignment
	buffer += c.alignCell(style)

	// right padding
	buffer += strings.Repeat(" ", style.PaddingRight)

	return buffer
}

func (c *Cell) alignCell(style *renderStyle) string {
	buffer := ""
	width := style.CellWidth(c.column)

	if c.colSpan > 1 {
		for i := 1; i < c.colSpan; i++ {
			w := style.CellWidth(c.column+i)
			if w == 0 {
				break
			}
			width += style.PaddingLeft + w + style.PaddingRight + len(style.BorderY)
		}
	}

	switch *c.alignment {

	default:
		buffer += c.formattedValue
		if l := width - c.Width(); l > 0 {
			buffer += strings.Repeat(" ", l)
		}

	case AlignLeft:
		buffer += c.formattedValue
		if l := width - c.Width(); l > 0 {
			buffer += strings.Repeat(" ", l)
		}

	case AlignRight:
		if l := width - c.Width(); l > 0 {
			buffer += strings.Repeat(" ", l)
		}
		buffer += c.formattedValue

	case AlignCenter:
		left, right :=  0, 0
		if l := width - c.Width(); l > 0 {
			lf := float64(l)
			left = int(math.Floor(lf / 2))
			right = int(math.Ceil(lf / 2))
		}
		buffer += strings.Repeat(" ", left)
		buffer += c.formattedValue
		buffer += strings.Repeat(" ", right)
	}

	return buffer
}

// Format the raw value as a string depending on the type
func renderValue(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case bool:
		return strconv.FormatBool(v.(bool))
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case uint64:
		return strconv.FormatUint(v.(uint64), 10)
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', 2, 64)
	}
	return ""
}
