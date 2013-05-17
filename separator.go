// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

import "strings"

type lineType int

const (
	// LINE_INNER *must* be the default
	LINE_INNER  lineType = iota
	LINE_TOP             // only descenders
	LINE_SUBTOP          // only descenders in the middle, but both at edges
	LINE_BOTTOM          // only ascenders
)

type Separator struct {
	where lineType
}

func (s *Separator) Render(style *renderStyle) string {
	// loop over getting dashes
	parts := []string{}
	for i := 0; i < style.columns; i++ {
		w := style.PaddingLeft + style.CellWidth(i) + style.PaddingRight
		parts = append(parts, strings.Repeat(style.BorderX, w))
	}

	switch s.where {
	case LINE_TOP:
		return style.BorderTopLeft + strings.Join(parts, style.BorderTop) + style.BorderTopRight
	case LINE_SUBTOP:
		return style.BorderLeft + strings.Join(parts, style.BorderTop) + style.BorderRight
	case LINE_BOTTOM:
		return style.BorderBottomLeft + strings.Join(parts, style.BorderBottom) + style.BorderBottomRight
	case LINE_INNER:
		return style.BorderLeft + strings.Join(parts, style.BorderI) + style.BorderRight
	}
	panic("not reached")
}
