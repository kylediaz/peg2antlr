package main

import (
	"testing"
)

type astTest struct {
	input  string
	output *node32
}

func TestPeg2Ast(t *testing.T) {
	tests := []astTest{
		{
			input:  "ruleName <- expression",
			output: nil,
		},
	}
	for _, test := range tests {
		fullInput := wrapWithPegTemplate(test.input)
		_, err := peg2ast(fullInput)
		if err != nil {
			t.Error("Error parsing input", test.input, ":", err)
		}
	}
}
