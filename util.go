package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var pegTemplate string = "package main\ntype Peg Peg {\n*tree.Tree\n}\n%s"

func wrapWithPegTemplate(input ...any) string {
	return fmt.Sprintf(pegTemplate, input...)
}

func writeBufferToFile(buffer bytes.Buffer, fileName string) error {
	err := os.WriteFile(fileName, buffer.Bytes(), 0644)
	return err
}

var tokenNameGenerateCount = 0
var tokenLabelGuesses = map[string]string{
	"[0-9]+": "INT",

	".": "DOT",
	",": "COMMA",
	"*": "STAR",
	"-": "DASH",

	"+":  "PLUS",
	"=":  "EQUAL",
	"&":  "AMPERSAND",
	"|":  "PIPE",
	"/":  "SLASH",
	"\\": "BACKSLASH",
	":":  "COLON",
	"::": "DOUBLE_COLON",
	";":  "SEMICOLON",
	"^":  "CARET",
	"$":  "DOLLAR",
	"?":  "QUESTION",

	"<":  "LT",
	">":  "GT",
	"<=": "LTE",
	">=": "GTE",
	"<>": "NEQ",

	"{": "LCURLY",
	"}": "RCURLY",
	"(": "LPAREN",
	")": "RPAREN",
	"[": "LBRACKET",
	"]": "RBRACKET",
}

var isAlpha = regexp.MustCompile("^[a-zA-Z]+$")

func generateTokenLabel(token string) string {
	if g, ok := tokenLabelGuesses[token]; ok {
		return g
	} else if isAlpha.MatchString(token) {
		return strings.ToUpper(token)
	} else {
		tokenNameGenerateCount += 1
		return fmt.Sprintf("TOKEN%d", tokenNameGenerateCount)
	}
}
