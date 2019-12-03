// Copyright 2012-2013 Apcera Inc. All rights reserved.

package termtables

import (
	"bytes"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/scylladb/termtables/term"
)

// MaxColumns represents the maximum number of columns that are available for
// display without wrapping around the right-hand side of the terminal window.
// At program initialization, the value will be automatically set according
// to available sources of information, including the $COLUMNS environment
// variable and, on Unix, tty information.
var MaxColumns = 80

// Element the interface that can draw a representation of the contents of a
// table cell.
type Element interface {
	Render(*renderStyle) string
}

type outputMode int

const (
	outputTerminal outputMode = iota
	outputMarkdown
	outputHTML
)

// Open question: should UTF-8 become an output mode?  It does require more
// tracking when resetting, if the locale-enabling had been used.

var outputsEnabled struct {
	UTF8       bool
	HTML       bool
	Markdown   bool
	titleStyle titleStyle
}

var defaultOutputMode outputMode = outputTerminal

// Table represents a terminal table.  The Style can be directly accessed
// and manipulated; all other access is via methods.
type Table struct {
	Style *TableStyle

	elements   []Element
	headers    []interface{}
	title      interface{}
	titleCell  *Cell
	outputMode outputMode
}

// EnableUTF8 will unconditionally enable using UTF-8 box-drawing characters
// for any tables created after this call, as the default style.
func EnableUTF8() {
	outputsEnabled.UTF8 = true
}

// SetModeHTML will control whether or not new tables generated will be in HTML
// mode by default; HTML-or-not takes precedence over options which control how
// a terminal output will be rendered, such as whether or not to use UTF8.
// This affects any tables created after this call.
func SetModeHTML(onoff bool) {
	outputsEnabled.HTML = onoff
	chooseDefaultOutput()
}

// SetModeMarkdown will control whether or not new tables generated will be
// in Markdown mode by default.  HTML-mode takes precedence.
func SetModeMarkdown(onoff bool) {
	outputsEnabled.Markdown = onoff
	chooseDefaultOutput()
}

var utfRe = regexp.MustCompile(`utf\-8|utf8|UTF\-8|UTF8`)

// EnableUTF8PerLocale will use current locale character map information to
// determine if UTF-8 is expected and, if so, is equivalent to EnableUTF8.
func EnableUTF8PerLocale() {
	locale := getLocale()
	if utfRe.MatchString(locale) {
		EnableUTF8()
	}
}

// getLocale returns the current locale name.
func getLocale() string {
	if runtime.GOOS == "windows" {
		// TODO: detect windows locale
		return "US-ASCII"
	}
	return unixLocale()
}

// unixLocale returns the locale by checking the $LC_ALL, $LC_CTYPE, and $LANG
// environment variables. If none of those are set, it returns "US-ASCII".
func unixLocale() string {
	for _, env := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if locale := os.Getenv(env); locale != "" {
			return locale
		}
	}
	return "US-ASCII"
}

// SetHTMLStyleTitle lets an HTML title output mode be chosen.
func SetHTMLStyleTitle(want titleStyle) {
	outputsEnabled.titleStyle = want
}

// chooseDefaultOutput sets defaultOutputMode based on priority
// choosing amongst the options which are enabled.  Pros: simpler
// encapsulation; cons: setting markdown doesn't disable HTML if
// HTML was previously enabled and was later disabled.
// This seems fairly reasonable.
func chooseDefaultOutput() {
	if outputsEnabled.HTML {
		defaultOutputMode = outputHTML
	} else if outputsEnabled.Markdown {
		defaultOutputMode = outputMarkdown
	} else {
		defaultOutputMode = outputTerminal
	}
}

func init() {
	// Do not enable UTF-8 per locale by default, breaks tests.
	sz, err := term.GetSize()
	if err == nil && sz.Columns != 0 {
		MaxColumns = sz.Columns
	}
}

// CreateTable creates an empty Table using defaults for style.
func CreateTable() *Table {
	t := &Table{elements: []Element{}, Style: DefaultStyle}
	if outputsEnabled.UTF8 {
		t.Style.setUtfBoxStyle()
	}
	if outputsEnabled.titleStyle != titleStyle(0) {
		t.Style.htmlRules.title = outputsEnabled.titleStyle
	}
	t.outputMode = defaultOutputMode
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
}

// AddHeaders supplies column headers for the table.
func (t *Table) AddHeaders(headers ...interface{}) {
	t.headers = append(t.headers, headers...)
}

// SetColumnAlignment changes the alignment for elements in a column of the table;
// alignments are stored with each cell, so cells added after a call to
// SetColumnAlignment will not pick up the change.  Columns are numbered from 1.
func (t *Table) SetColumnAlignment(align TableAlignment, columns ...int) {
	for i := range t.elements {
		row, ok := t.elements[i].(*Row)
		if !ok {
			continue
		}
		for _, column := range columns {
			if column < 0 || column >= len(row.cells) {
				continue
			}
			row.cells[column-1].alignment = &align
		}
	}
}

// UTF8Box sets the table style to use UTF-8 box-drawing characters,
// overriding all relevant style elements at the time of the call.
func (t *Table) UTF8Box() {
	t.Style.setUtfBoxStyle()
}

// SetModeHTML switches this table to be in HTML when rendered; the
// default depends upon whether the package function SetModeHTML() has been
// called, and with what value.  This method forces the feature on for this
// table.  Turning off involves choosing a different mode, per-table.
func (t *Table) SetModeHTML() {
	t.outputMode = outputHTML
}

// SetModeMarkdown switches this table to be in Markdown mode
func (t *Table) SetModeMarkdown() {
	t.outputMode = outputMarkdown
}

// SetModeTerminal switches this table to be in terminal mode.
func (t *Table) SetModeTerminal() {
	t.outputMode = outputTerminal
}

// SetHTMLStyleTitle lets an HTML output mode be chosen; we should rework this
// into a more generic and extensible API as we clean up termtables.
func (t *Table) SetHTMLStyleTitle(want titleStyle) {
	t.Style.htmlRules.title = want
}

// Render returns a string representation of a fully rendered table, drawn
// out for display, with embedded newlines.  If this table is in HTML mode,
// then this is equivalent to RenderHTML().
func (t *Table) Render() string {
	// Elements is already populated with row data.
	switch t.outputMode {
	case outputTerminal:
		return t.renderTerminal()
	case outputMarkdown:
		return t.renderMarkdown()
	case outputHTML:
		return t.RenderHTML()
	default:
		panic("unknown output mode set")
	}
}

// renderTerminal returns a string representation of a fully rendered table,
// drawn out for display, with embedded newlines.
func (t *Table) renderTerminal() string {
	// Use a placeholder rather than adding titles/headers to the tables
	// elements or else successive calls will compound them.
	tt := t.clone()

	// Initial top line.
	if !tt.Style.SkipBorder {
		if tt.title != nil && tt.headers == nil {
			tt.elements = append([]Element{&Separator{where: LINE_SUBTOP}}, tt.elements...)
		} else if tt.title == nil && tt.headers == nil {
			tt.elements = append([]Element{&Separator{where: LINE_TOP}}, tt.elements...)
		} else {
			tt.elements = append([]Element{&Separator{where: LINE_INNER}}, tt.elements...)
		}
	}

	// If we have headers, include them.
	if tt.headers != nil {
		ne := make([]Element, 2)
		ne[1] = CreateRow(tt.headers)
		if tt.title != nil {
			ne[0] = &Separator{where: LINE_SUBTOP}
		} else {
			ne[0] = &Separator{where: LINE_TOP}
		}
		tt.elements = append(ne, tt.elements...)
	}

	// If we have a title, write it.
	if tt.title != nil {
		// Match changes to this into renderMarkdown too.
		tt.titleCell = CreateCell(tt.title, &CellStyle{Alignment: AlignCenter, ColSpan: 999})
		ne := []Element{
			&StraightSeparator{where: LINE_TOP},
			CreateRow([]interface{}{tt.titleCell}),
		}
		tt.elements = append(ne, tt.elements...)
	}

	// Create a new table from the
	// generate the runtime style. Must include all cells being printed.
	style := createRenderStyle(tt)

	// Loop over the elements and render them.
	b := bytes.NewBuffer(nil)
	for _, e := range tt.elements {
		b.WriteString(e.Render(style))
		b.WriteString("\n")
	}

	// Add bottom line.
	if !style.SkipBorder {
		b.WriteString((&Separator{where: LINE_BOTTOM}).Render(style) + "\n")
	}

	return b.String()
}

// renderMarkdown returns a string representation of a table in Markdown
// markup format using GitHub Flavored Markdown's notation (since tables
// are not in the core Markdown spec).
func (t *Table) renderMarkdown() string {
	// We need ASCII drawing characters; we need a line after the header;
	// *do* need a header!  Do not need to markdown-escape contents of
	// tables as markdown is ignored in there.  Do need to do _something_
	// with a '|' character shown as a member of a table.

	t.Style.setAsciiBoxStyle()

	firstLines := make([]Element, 0, 2)

	if t.headers == nil {
		initial := createRenderStyle(t)
		if initial.columns > 1 {
			row := CreateRow([]interface{}{})
			for i := 0; i < initial.columns; i++ {
				row.AddCell(CreateCell(i+1, &CellStyle{}))
			}
		}
	}

	firstLines = append(firstLines, CreateRow(t.headers))
	// This is a dummy line, swapped out below.
	firstLines = append(firstLines, firstLines[0])
	t.elements = append(firstLines, t.elements...)
	// Generate the runtime style.
	style := createRenderStyle(t)
	// We know that the second line is a dummy, we can replace it.
	mdRow := CreateRow([]interface{}{})
	for i := 0; i < style.columns; i++ {
		mdRow.AddCell(CreateCell(strings.Repeat("-", style.cellWidths[i]), &CellStyle{}))
	}
	t.elements[1] = mdRow

	b := bytes.NewBuffer(nil)
	// Comes after style is generated, which must come after all width-affecting
	// changes are in.
	if t.title != nil {
		// Markdown doesn't support titles or column spanning; we _should_
		// escape the title, but doing that to handle all possible forms of
		// markup would require a heavy dependency, so we punt.
		b.WriteString("Table: ")
		b.WriteString(strings.TrimSpace(CreateCell(t.title, &CellStyle{}).Render(style)))
		b.WriteString("\n\n")
	}

	// Loop over the elements and render them.
	for _, e := range t.elements {
		b.WriteString(e.Render(style))
		b.WriteString("\n")
	}

	return b.String()
}

// clone returns a copy of the table with the underlying slices being copied;
// the references to the Elements/cells are left as shallow copies.
func (t *Table) clone() *Table {
	tt := &Table{outputMode: t.outputMode, Style: t.Style, title: t.title}
	if t.headers != nil {
		tt.headers = make([]interface{}, len(t.headers))
		copy(tt.headers, t.headers)
	}
	if t.elements != nil {
		tt.elements = make([]Element, len(t.elements))
		copy(tt.elements, t.elements)
	}
	return tt
}
