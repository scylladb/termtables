# Termtables

A [Go](http://golang.org) port of the Ruby library [terminal-tables][https://github.com/visionmedia/terminal-table] for
fast and simple ASCII table generation.

## Installation

```bash
go get github.com/apcera/termtables
```

## Go Style Documentation

[http://go.pkgdoc.org/github.com/apcera/termtable](http://go.pkgdoc.org/github.com/apcera/termtable)

## Basic Usage

```go
package main

import (
  "fmt"
  "github.com/apcera/termtables"
)

func main() {
  table := termtables.CreateTable()

  table.AddHeaders("Name", "Age")
  table.AddRow("John", "30")
  table.AddRow("Sam", 18)
  table.AddRow("Julie", 20.14)

  fmt.Println(table.Render())
}
```

```
+-------+-------+
| Name  | Age   |
+-------+-------+
| John  | 30    |
| Sam   | 18    |
| Julie | 20.14 |
+-------+-------+
```
