package storeaccesslog_test

import "regexp"

func normalize(in string) string { return regexp.QuoteMeta(in) }
