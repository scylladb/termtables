package termtables

import (
	"testing"
)

func TestCellRenderString(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := createCell(0, "foobar", nil)

	output := cell.Render(style)
	if output != "foobar" {
		t.Fatal("Unexpected output:", output)
	}
}

func TestCellRenderBool(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := createCell(0, true, nil)

	output := cell.Render(style)
	if output != "true" {
		t.Fatal("Unexpected output:", output)
	}
}

func TestCellRenderInteger(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := createCell(0, 12345, nil)

	output := cell.Render(style)
	if output != "12345" {
		t.Fatal("Unexpected output:", output)
	}
}

func TestCellRenderFloat(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := createCell(0, 12.345, nil)

	output := cell.Render(style)
	if output != "12.35" {
		t.Fatal("Unexpected output:", output)
	}
}

func TestCellRenderPadding(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{ PaddingLeft: 3, PaddingRight: 4 }, cellWidths: map[int]int{} }

	cell := createCell(0, "foobar", nil)

	output := cell.Render(style)
	if output != "   foobar    " {
		t.Fatal("Unexpected output:", output)
	}
}

type foo struct {
	v string
}
func (f *foo) String() string {
	return f.v
}

func TestCellRenderStringerStruct(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := createCell(0, &foo{v:"bar"}, nil)

	output := cell.Render(style)
	if output != "bar" {
		t.Fatal("Unexpected output:", output)
	}
}

type fooString string

func TestCellRenderGeneric(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := createCell(0, fooString("baz"), nil)

	output := cell.Render(style)
	if output != "baz" {
		t.Fatal("Unexpected output:", output)
	}
}
