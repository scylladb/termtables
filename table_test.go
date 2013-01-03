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

func TestCreateTable(t *testing.T) {
	expected := "+-----------+-------+\n" +
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

	output := table.Render()
	if output != expected {
		t.Fatal(DisplayFailedOutput(output, expected))
	}
}

func TestTableWithHeader(t *testing.T) {
	expected := "+-------------------+\n" +
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

	output := table.Render()
	if output != expected {
		t.Fatal(DisplayFailedOutput(output, expected))
	}
}

func TestTableWithNoHeaders(t *testing.T) {
	expected := "+-----------+------+\n" +
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

	output := table.Render()
	if output != expected {
		t.Fatal(DisplayFailedOutput(output, expected))
	}
}

func TestTableUnicodeWidths(t *testing.T) {
	expected := "+-----------+------+\n" +
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

	output := table.Render()
	if output != expected {
		t.Fatal(DisplayFailedOutput(output, expected))
	}
}
