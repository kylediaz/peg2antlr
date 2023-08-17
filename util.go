package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func writeBufferToFile(buffer bytes.Buffer, fileName string) error {
	err := os.WriteFile(fileName, buffer.Bytes(), 0644)
	return err
}

// Reverses linked list
func reverse(node *node32) *node32 {
	var head *node32
	for node != nil {
		next := node.next
		node.next = head
		head = node
		node = next
	}
	return head
}

func getRawNodeValue(node *node32, buffer string) string {
	return buffer[node.begin:node.end]
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

// DFS prefix
func traverse(root *node32, f func(node *node32)) {
	for root != nil {
		f(root)
		if root.up != nil {
			traverse(root.up, f)
		}
		root = root.next
	}
}

// Returns all nodes matching the predicate
func collect(root *node32, predicate func(node *node32) bool) []*node32 {
	output := make([]*node32, 0, 10)
	traverse(root, func(node *node32) {
		if predicate(node) {
			output = append(output, node)
		}
	})
	return output
}

func remove(root *node32, ruleType pegRule) *node32 {
	if root == nil {
		return nil
	} else if root.pegRule == ruleType {
		return nil
	} else {
		root.up = remove(root.up, ruleType)
		root.next = remove(root.next, ruleType)
		return root
	}
}
