// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

type Element interface {
	Render(*renderStyle) string
}

type Table struct {
	Style *TableStyle

	elements []Element
	headers *[]interface{}
	title   *interface{}
}

func CreateTable() *Table {
	return &Table{ elements: []Element{}, Style: DefaultStyle }
}

func (t *Table) AddSeparator() {
	t.elements = append(t.elements, &Separator{})
}

func (t *Table) AddRow(items ...interface{}) {
	t.elements = append(t.elements, CreateRow(items))
}

func (t *Table) AddTitle(title interface{}) {
	t.title = &title
}

func (t *Table) AddHeaders(headers ...interface{}) {
	t.headers = &headers
}

func (t *Table) Render() (buffer string) {
	// used generically to add a separator
	separator := Separator{}

	// initial top line
	if !t.Style.SkipBorder {
		ne := []Element { &separator }
		ne = append(ne, t.elements...)
		t.elements = ne
	}

	// if we have headers, include them
	if t.headers != nil {
		// FIXME: there must be a better way to prepend an array in Go than this
		ne := []Element { &separator, CreateRow(*t.headers) }
		ne = append(ne, t.elements...)
		t.elements = ne
	}

	// if we have a title, write them
	if t.title != nil {
		ne := []Element { &StraightSeparator{}, CreateRow([]interface{}{CreateCell(*t.title, &CellStyle{Alignment: AlignCenter, ColSpan: 999})}) }
		ne = append(ne, t.elements...)
		t.elements = ne
	}

	// generate the runtime style
	style := createRenderStyle(t)

	// loop over the elements and render them
	for _, e := range t.elements {
		buffer += e.Render(style) + "\n"
	}

	// add bottom line
	if !style.SkipBorder {
		buffer += separator.Render(style) + "\n"
	}

	return buffer
}
