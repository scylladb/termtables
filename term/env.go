// Copyright 2013 Apcera Inc. All rights reserved.

package term

import (
	"os"
	"strconv"
)

func GetEnvWindowSize() (*Size, bool) {
	lines := os.Getenv("LINES")
	columns := os.Getenv("COLUMNS")
	if lines == "" && columns == "" {
		return nil, false
	}

	nLines := 0
	nColumns := 0
	var err error
	if lines != "" {
		nLines, err = strconv.Atoi(lines)
		if err != nil {
			return nil, false
		}
	}
	if columns != "" {
		nColumns, err = strconv.Atoi(columns)
		if err != nil {
			return nil, false
		}
	}

	return &Size{
		Lines:   nLines,
		Columns: nColumns,
	}, true
}
