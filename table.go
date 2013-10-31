// Copyright 2012-2013 Apcera Inc. All rights reserved.

package termtables

import (
	"strings"

	"github.com/apcera/termtables/locale"
	"github.com/apcera/termtables/term"
)

var useUTF8ByDefault = false
var useHTMLByDefault = false

// MaxColumns represents the maximum number of columns that are available for
// display without wrapping around the right-hand side of the terminal window.
// At program initialization, the value will be automatically set according
// to available sources of information, including the $COLUMNS environment
// variable and, on Unix, tty information.
var MaxColumns = 80

// An Element is a drawn representation of the contents of a table cell.
type Element interface {
	Render(*renderStyle) string
}

type outputMode int

const (
	outputTerminal outputMode = iota
	outputHTML
)

// Table represents a terminal table.  The Style can be directly accessed
// and manipulated; all other access is via methods.
type Table struct {
	Style *TableStyle

	elements   []Element
	minWidth   int
	headers    []interface{}
	title      interface{}
	outputMode outputMode
}

// EnableUTF8 will unconditionally enable using UTF-8 box-drawing characters
// for any tables created after this call, as the default style.
func EnableUTF8() {
	useUTF8ByDefault = true
}

// SetModeHTML will control whether or not new tables generated will be in HTML
// mode by default; HTML-or-not takes precedence over options which control how
// a terminal output will be rendered, such as whether or not to use UTF8.
// This affects any tables created after this call.
func SetModeHTML(onoff bool) {
	useHTMLByDefault = onoff
}

// EnableUTF8PerLocale will use current locale character map information to
// determine if UTF-8 is expected and, if so, is equivalent to EnableUTF8.
func EnableUTF8PerLocale() {
	charmap := locale.GetCharmap()
	if strings.EqualFold(charmap, "UTF-8") {
		useUTF8ByDefault = true
	}
}

func init() {
	// do not enable UTF-8 per locale by default, breaks tests
	sz, err := term.GetSize()
	if err == nil && sz.Columns != 0 {
		MaxColumns = sz.Columns
	}
}

// CreateTable creates an empty Table using defaults for style.
func CreateTable() *Table {
	t := &Table{elements: []Element{}, Style: DefaultStyle}
	if useUTF8ByDefault {
		t.Style.setUtfBoxStyle()
	}
	if useHTMLByDefault {
		t.outputMode = outputHTML
	}
	return t
}

// AddSeparator adds a line to the table content, where the line
// consists of separator characters.
func (t *Table) AddSeparator() {
	t.elements = append(t.elements, &Separator{})
}

// AddRow adds the supplied items as cells in one row of the table.
func (t *Table) AddRow(items ...interface{}) *Row {
	row := CreateRow(items)
	t.elements = append(t.elements, row)
	return row
}

// AddTitle supplies a table title, which if present will be rendered as
// one cell across the width of the table, as the first row.
func (t *Table) AddTitle(title interface{}) {
	t.title = title

	t.minWidth = len(renderValue(title))
}

// AddHeaders supplies column headers for the table.
func (t *Table) AddHeaders(headers ...interface{}) {
	t.headers = headers[:]
}

// UTF8Box sets the table style to use UTF-8 box-drawing characters,
// overriding all relevant style elements at the time of the call.
func (t *Table) UTF8Box() {
	t.Style.setUtfBoxStyle()
}

// SetModeHTML will control whether or not this table will be emitted as
// HTML when rendered; the default depends upon whether the package function
// SetModeHTML() has been called, and with what value.  This method controls
// the same functionality, but on a per-table basis.
func (t *Table) SetModeHTML(onoff bool) {
	if onoff {
		t.outputMode = outputHTML
	} else {
		t.outputMode = outputTerminal
	}
}

// Render returns a string representation of a fully rendered table, drawn
// out for display, with embedded newlines.  If this table is in HTML mode,
// then this is equivalent to RenderHTML().
func (t *Table) Render() (buffer string) {
	// elements is already populated with row data
	switch t.outputMode {
	case outputTerminal:
		return t.renderTerminal()
	case outputHTML:
		return t.RenderHTML()
	default:
		panic("unknown output mode set")
	}
}

// renderTerminal returns a string representation of a fully rendered table,
// drawn out for display, with embedded newlines.
func (t *Table) renderTerminal() (buffer string) {
	// initial top line
	if !t.Style.SkipBorder {
		if t.title != nil && t.headers == nil {
			t.elements = append([]Element{&Separator{where: LINE_SUBTOP}}, t.elements...)
		} else if t.title == nil && t.headers == nil {
			t.elements = append([]Element{&Separator{where: LINE_TOP}}, t.elements...)
		} else {
			t.elements = append([]Element{&Separator{where: LINE_INNER}}, t.elements...)
		}
	}

	// if we have headers, include them
	if t.headers != nil {
		ne := make([]Element, 2)
		ne[1] = CreateRow(t.headers)
		if t.title != nil {
			ne[0] = &Separator{where: LINE_SUBTOP}
		} else {
			ne[0] = &Separator{where: LINE_TOP}
		}
		t.elements = append(ne, t.elements...)
	}

	// if we have a title, write them
	if t.title != nil {
		ne := []Element{
			&StraightSeparator{where: LINE_TOP},
			CreateRow([]interface{}{CreateCell(t.title, &CellStyle{Alignment: AlignCenter, ColSpan: 999})}),
		}
		t.elements = append(ne, t.elements...)
	}

	// generate the runtime style
	style := createRenderStyle(t)

	// loop over the elements and render them
	for _, e := range t.elements {
		buffer += e.Render(style) + "\n"
	}

	// add bottom line
	if !style.SkipBorder {
		buffer += (&Separator{where: LINE_BOTTOM}).Render(style) + "\n"
	}

	return buffer
}
