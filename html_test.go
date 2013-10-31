// Copyright 2013 Apcera Inc. All rights reserved.

package termtables

import (
	"testing"
)

func TestCreateTableHTML(t *testing.T) {
	expected := "<table>\n" +
		"<tr><th>Name     </th><th>Value</th></tr>\n" +
		"<tr><td>hey      </td><td>you </td></tr>\n" +
		"<tr><td>ken      </td><td>1234</td></tr>\n" +
		"<tr><td>derek    </td><td>3.14</td></tr>\n" +
		"<tr><td>derek too</td><td>3.15</td></tr>\n" +
		"</table>\n"

	table := CreateTable()
	table.SetModeHTML()

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

func TestTableWithHeaderHTML(t *testing.T) {
	expected := "<table>\n" +
		"<caption>Example</caption>\n" +
		"<tr><th>Name     </th><th>Value</th></tr>\n" +
		"<tr><td>hey      </td><td>you </td></tr>\n" +
		"<tr><td>ken      </td><td>1234</td></tr>\n" +
		"<tr><td>derek    </td><td>3.14</td></tr>\n" +
		"<tr><td>derek too</td><td>3.15</td></tr>\n" +
		"</table>\n"

	table := CreateTable()
	table.SetModeHTML()

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

func TestTableTitleWidthAdjustsHTML(t *testing.T) {
	expected := "<table>\n" +
		"<caption>Example My Foo Bar&#39;d Test</caption>\n" +
		"<tr><th>Name     </th><th>Value              </th></tr>\n" +
		"<tr><td>hey      </td><td>you                </td></tr>\n" +
		"<tr><td>ken      </td><td>1234               </td></tr>\n" +
		"<tr><td>derek    </td><td>3.14               </td></tr>\n" +
		"<tr><td>derek too</td><td>3.15               </td></tr>\n" +
		"</table>\n"

	table := CreateTable()
	table.SetModeHTML()

	table.AddTitle("Example My Foo Bar'd Test")
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

func TestTableWithNoHeadersHTML(t *testing.T) {
	expected := "<table>\n" +
		"<tr><td>hey      </td><td>you </td></tr>\n" +
		"<tr><td>ken      </td><td>1234</td></tr>\n" +
		"<tr><td>derek    </td><td>3.14</td></tr>\n" +
		"<tr><td>derek too</td><td>3.15</td></tr>\n" +
		"</table>\n"

	table := CreateTable()
	table.SetModeHTML()

	table.AddRow("hey", "you")
	table.AddRow("ken", 1234)
	table.AddRow("derek", 3.14)
	table.AddRow("derek too", 3.1456788)

	output := table.Render()
	if output != expected {
		t.Fatal(DisplayFailedOutput(output, expected))
	}
}

func TestTableUnicodeWidthsHTML(t *testing.T) {
	expected := "<table>\n" +
		"<tr><th>Name     </th><th>Cost</th></tr>\n" +
		"<tr><td>Currency </td><td>¤10</td></tr>\n" +
		"<tr><td>US Dollar</td><td>$30</td></tr>\n" +
		"<tr><td>Euro     </td><td>€27</td></tr>\n" +
		"<tr><td>Thai     </td><td>฿70</td></tr>\n" +
		"</table>\n"

	table := CreateTable()
	table.SetModeHTML()
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
