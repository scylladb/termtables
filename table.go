// Copyright 2012 Apcera Inc. All rights reserved.

package termtables

type Element interface {
	Render(*renderStyle) string
}

type Table struct {
	Style *TableStyle

	elements []Element
	headers *[]interface{}
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

func (t *Table) AddHeaders(headers ...interface{}) {
	t.headers = &headers
}

func (t *Table) Render() (buffer string) {
	// used generically to add a separator
	separator := Separator{}

	// if we have headers, include them
	if t.headers != nil {
		// FIXME: there must be a better way to prepend an array in Go than this
		ne := []Element { CreateRow(*t.headers), &separator }
		ne = append(ne, t.elements...)
		t.elements = ne
	}

	// generate the runtime style
	style := createRenderStyle(t)

	// initial top line
	buffer += separator.Render(style) + "\n"

	// loop over the elements and render them
	for _, e := range t.elements {
		buffer += e.Render(style) + "\n"
	}

	// add bottom line
	buffer += separator.Render(style) + "\n"

	return buffer
}
