package termtables

import (
	"testing"
)

func TestCellRenderString(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := CreateCell(0, "foobar")

	output := cell.Render(style)
	if output != "foobar" {
		t.Fatal("Unexpected output:", output)
	}
}

func TestCellRenderBool(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := CreateCell(0, true)

	output := cell.Render(style)
	if output != "true" {
		t.Fatal("Unexpected output:", output)
	}
}

func TestCellRenderInteger(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := CreateCell(0, 12345)

	output := cell.Render(style)
	if output != "12345" {
		t.Fatal("Unexpected output:", output)
	}
}

func TestCellRenderFloat(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{}, cellWidths: map[int]int{} }
	cell := CreateCell(0, 12.345)

	output := cell.Render(style)
	if output != "12.35" {
		t.Fatal("Unexpected output:", output)
	}
}

func TestCellRenderPadding(t *testing.T) {
	style := &renderStyle{ TableStyle: TableStyle{ PaddingLeft: 3, PaddingRight: 4 }, cellWidths: map[int]int{} }

	cell := CreateCell(0, "foobar")

	output := cell.Render(style)
	if output != "   foobar    " {
		t.Fatal("Unexpected output:", output)
	}
}
