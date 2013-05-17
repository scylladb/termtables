// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

import "strings"

type Separator struct{}

func (s *Separator) Render(style *renderStyle) string {
	// loop over getting dashes
	parts := []string{}
	for i := 0; i < style.columns; i++ {
		w := style.PaddingLeft + style.CellWidth(i) + style.PaddingRight
		parts = append(parts, strings.Repeat(style.BorderX, w))
	}

	return style.BorderI + strings.Join(parts, style.BorderI) + style.BorderI
}
