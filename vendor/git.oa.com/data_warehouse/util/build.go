package util

/*
#include <stdio.h>
#include <string.h>

void get_date_from_macro(const char* date_str, int* year, int* month, int*day)
{
	static const char month_names[] = "JanFebMarAprMayJunJulAugSepOctNovDec";
	char month_str[5];

	sscanf(date_str, "%s %d %d", month_str, day, year);
	*month = (strstr(month_names, month_str) - month_names) / 3 + 1;

	return;
}

const char* build_time(void)
{
	static char sz_build_time[20];
	int year, month, day;

	memset(sz_build_time, 0, 20);

	get_date_from_macro(__DATE__, &year, &month, &day);
	snprintf(sz_build_time, 20, "%04d-%02d-%02d %s", year, month, day, __TIME__);

	return sz_build_time;
}
*/
import "C"

func BuildTime() string {
	return C.GoString(C.build_time())
}
