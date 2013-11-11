package termtables

import (
	"testing"
)

func DisplayFailedOutput(actual, expected string) string {
	return "Output didn't match expected\n\n" +
		"Actual:\n\n" +
		actual + "\n" +
		"Expected:\n\n" +
		expected
}

func checkRendersTo(t *testing.T, table *Table, expected string) {
	output := table.Render()
	if output != expected {
		t.Fatal(DisplayFailedOutput(output, expected))
	}
}

func TestCreateTable(t *testing.T) {
	expected := "" +
		"+-----------+-------+\n" +
		"| Name      | Value |\n" +
		"+-----------+-------+\n" +
		"| hey       | you   |\n" +
		"| ken       | 1234  |\n" +
		"| derek     | 3.14  |\n" +
		"| derek too | 3.15  |\n" +
		"+-----------+-------+\n"

	table := CreateTable()

	table.AddHeaders("Name", "Value")
	table.AddRow("hey", "you")
	table.AddRow("ken", 1234)
	table.AddRow("derek", 3.14)
	table.AddRow("derek too", 3.1456788)

	checkRendersTo(t, table, expected)
}

func TestStyleResets(t *testing.T) {
	expected := "" +
		"+-----------+-------+\n" +
		"| Name      | Value |\n" +
		"+-----------+-------+\n" +
		"| hey       | you   |\n" +
		"| ken       | 1234  |\n" +
		"| derek     | 3.14  |\n" +
		"| derek too | 3.15  |\n" +
		"+-----------+-------+\n"

	table := CreateTable()
	table.UTF8Box()
	table.Style.setAsciiBoxStyle()

	table.AddHeaders("Name", "Value")
	table.AddRow("hey", "you")
	table.AddRow("ken", 1234)
	table.AddRow("derek", 3.14)
	table.AddRow("derek too", 3.1456788)

	checkRendersTo(t, table, expected)
}

func TestTableWithHeader(t *testing.T) {
	expected := "" +
		"+-------------------+\n" +
		"|      Example      |\n" +
		"+-----------+-------+\n" +
		"| Name      | Value |\n" +
		"+-----------+-------+\n" +
		"| hey       | you   |\n" +
		"| ken       | 1234  |\n" +
		"| derek     | 3.14  |\n" +
		"| derek too | 3.15  |\n" +
		"+-----------+-------+\n"

	table := CreateTable()

	table.AddTitle("Example")
	table.AddHeaders("Name", "Value")
	table.AddRow("hey", "you")
	table.AddRow("ken", 1234)
	table.AddRow("derek", 3.14)
	table.AddRow("derek too", 3.1456788)

	checkRendersTo(t, table, expected)
}

func TestTableTitleWidthAdjusts(t *testing.T) {
	expected := "" +
		"+--------------------------------+\n" +
		"|   Example My Foo Bar'd Test    |\n" +
		"+-----------+--------------------+\n" +
		"| Name      | Value              |\n" +
		"+-----------+--------------------+\n" +
		"| hey       | you                |\n" +
		"| ken       | 1234               |\n" +
		"| derek     | 3.14               |\n" +
		"| derek too | 3.15               |\n" +
		"+-----------+--------------------+\n"

	table := CreateTable()

	table.AddTitle("Example My Foo Bar'd Test")
	table.AddHeaders("Name", "Value")
	table.AddRow("hey", "you")
	table.AddRow("ken", 1234)
	table.AddRow("derek", 3.14)
	table.AddRow("derek too", 3.1456788)

	checkRendersTo(t, table, expected)
}

func TestTableHeaderWidthAdjusts(t *testing.T) {
	expected := "" +
		"+---------------+---------------------+\n" +
		"| Slightly Long | More than 2 columns |\n" +
		"+---------------+---------------------+\n" +
		"| a             | b                   |\n" +
		"+---------------+---------------------+\n"

	table := CreateTable()

	table.AddHeaders("Slightly Long", "More than 2 columns")
	table.AddRow("a", "b")

	checkRendersTo(t, table, expected)
}

func TestTableWithNoHeaders(t *testing.T) {
	expected := "" +
		"+-----------+------+\n" +
		"| hey       | you  |\n" +
		"| ken       | 1234 |\n" +
		"| derek     | 3.14 |\n" +
		"| derek too | 3.15 |\n" +
		"+-----------+------+\n"

	table := CreateTable()

	table.AddRow("hey", "you")
	table.AddRow("ken", 1234)
	table.AddRow("derek", 3.14)
	table.AddRow("derek too", 3.1456788)

	checkRendersTo(t, table, expected)
}

func TestTableUnicodeWidths(t *testing.T) {
	expected := "" +
		"+-----------+------+\n" +
		"| Name      | Cost |\n" +
		"+-----------+------+\n" +
		"| Currency  | ¤10  |\n" +
		"| US Dollar | $30  |\n" +
		"| Euro      | €27  |\n" +
		"| Thai      | ฿70  |\n" +
		"+-----------+------+\n"

	table := CreateTable()
	table.AddHeaders("Name", "Cost")
	table.AddRow("Currency", "¤10")
	table.AddRow("US Dollar", "$30")
	table.AddRow("Euro", "€27")
	table.AddRow("Thai", "฿70")

	checkRendersTo(t, table, expected)
}

func TestTableInUTF8(t *testing.T) {
	expected := "" +
		"╭───────────────────╮\n" +
		"│      Example      │\n" +
		"├───────────┬───────┤\n" +
		"│ Name      │ Value │\n" +
		"├───────────┼───────┤\n" +
		"│ hey       │ you   │\n" +
		"│ ken       │ 1234  │\n" +
		"│ derek     │ 3.14  │\n" +
		"│ derek too │ 3.15  │\n" +
		"╰───────────┴───────╯\n"

	table := CreateTable()
	table.UTF8Box()

	table.AddTitle("Example")
	table.AddHeaders("Name", "Value")
	table.AddRow("hey", "you")
	table.AddRow("ken", 1234)
	table.AddRow("derek", 3.14)
	table.AddRow("derek too", 3.1456788)

	checkRendersTo(t, table, expected)
}

func TestTableInMarkdown(t *testing.T) {
	expected := "" +
		"Table: Example\n\n" +
		"| Name  | Value |\n" +
		"| ----- | ----- |\n" +
		"| hey   | you   |\n" +
		"| a &#x7c; b | esc   |\n"

	table := CreateTable()
	table.SetModeMarkdown()

	table.AddTitle("Example")
	table.AddHeaders("Name", "Value")
	table.AddRow("hey", "you")
	table.AddRow("a | b", "esc")

	checkRendersTo(t, table, expected)
}

func TestTitleUnicodeWidths(t *testing.T) {
	expected := "" +
		"+-------+\n" +
		"| ← 5 → |\n" +
		"+---+---+\n" +
		"| a | b |\n" +
		"| c | d |\n" +
		"| e | 3 |\n" +
		"+---+---+\n"

	// minimum width for a table of two columns is 9 characters, given
	// one space of padding, and non-empty tables.

	table := CreateTable()

	// We have 4 characters down for left and right columns and padding, so
	// a width of 5 for us should match the minimum per the columns

	// 5 characters; each arrow is three octets in UTF-8, giving 9 bytes
	// so, same in character-count-width, longer in bytes
	table.AddTitle("← 5 →")

	// a single character per cell, here; use ASCII characters
	table.AddRow("a", "b")
	table.AddRow("c", "d")
	table.AddRow("e", 3)

	checkRendersTo(t, table, expected)
}
