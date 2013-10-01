// Copyright 2013 Apcera Inc. All rights reserved.

package locale

/*
#include <langinfo.h>
#include <locale.h>

void
init_go_locale(void)
{
	(void) setlocale(LC_ALL, "");
}

const char *
get_charmap(void)
{
	return nl_langinfo(CODESET);
}

*/
import "C"

func init() {
	C.init_go_locale()
}

func GetCharmap() string {
	return C.GoString(C.get_charmap())
}
