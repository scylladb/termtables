// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

import "strings"

type StraightSeparator struct{}

func (s *StraightSeparator) Render(style *renderStyle) string {
	// loop over getting dashes
	width := 0
	for i := 0; i < style.columns; i++ {
		width += style.PaddingLeft + style.CellWidth(i) + style.PaddingRight + len(style.BorderI)
	}

	return style.BorderI + strings.Repeat(style.BorderX, width-1) + style.BorderI
}
