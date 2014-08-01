// Copyright 2014 Apcera Inc. All rights reserved.

// +build windows,!cgo

package locale

// FIXME: we need someone who understands Windows code-pages to make sure that
// we can use advanced characters before we move away from just saying
// "err, ASCII, *cough*".
func GetCharmap() string {
	return "US-ASCII"
}
