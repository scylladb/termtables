// Copyright 2013 Apcera Inc. All rights reserved.

package term

import (
	"os"
)

type Size struct {
	Lines   int
	Columns int
}

// Get the terminal window size.
// We prefer environ $LINES/$COLUMNS, then fall back to tty-held information.
// We do not support use of termcap/terminfo to derive size information.
func GetSize() (*Size, error) {
	envSize, envOk := GetEnvWindowSize()
	if envOk && envSize.Lines != 0 && envSize.Columns != 0 {
		return envSize, nil
	}

	fh, err := os.Open("/dev/tty")
	if err != nil {
		// no tty, no point continuing; we only let the environ
		// avoid an error in this case because if someone has faked
		// up an environ with LINES/COLUMNS _both_ set, we should let
		// them
		return nil, err
	}

	size, err := GetTerminalWindowSize(fh)
	if err != nil {
		if envOk {
			return envSize, nil
		}
		return nil, err
	}
	if !envOk {
		return size, err
	}

	if envSize.Lines == 0 {
		envSize.Lines = size.Lines
	}
	if envSize.Columns == 0 {
		envSize.Columns = size.Columns
	}
	return envSize, nil
}
