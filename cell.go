// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

import (
	"strconv"
	"strings"
)

type Cell struct {
	column int
	formattedValue string
}

func CreateCell(column int, v interface{}) *Cell {
	return &Cell{ column: column, formattedValue: renderValue(v) }
}

func (c *Cell) Width() int {
	return len(c.formattedValue)
}

func (c *Cell) Render(style *renderStyle) (buffer string) {
	// left padding
	buffer += strings.Repeat(" ", style.PaddingLeft)

	// FIXME: needs alignment handling

	// append the main value
	buffer += c.formattedValue

	// add padding
	if style.CellWidth(c.column) - c.Width() > 0 {
		buffer += strings.Repeat(" ", style.CellWidth(c.column) - c.Width())
	}

	// right padding
	buffer += strings.Repeat(" ", style.PaddingRight)

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
