package utils

import "regexp"

func splitByEmptyNewline(str, reg string) []string {

	return regexp.MustCompile(reg).Split(str, -1)

}
