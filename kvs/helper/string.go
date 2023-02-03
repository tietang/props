package helper

import "strings"

func Split(s string) []string {
	words := SplitAny(s, "/\\<>《》. -_()（）【】[]·─┅～—ˉ＿~﹣-、。\n〔〕\t\r")
	return words
}

func SplitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}
