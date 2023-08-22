package main

import (
	"bytes"
	"testing"
)

type testcase struct {
	input, output string
}

func TestPeg2AntlrUnit(t *testing.T) {
	tests := []testcase{
		{
			input:  "ruleName <- expression",
			output: "ruleName: expression;\n",
		},
		{
			input:  "ruleName <- e1 / e2",
			output: "ruleName: e1 | e2;\n",
		},
	}
	for _, test := range tests {
		fullInput := wrapWithPegTemplate(test.input)
		var dummy, parserOutputBuffer bytes.Buffer
		err := Peg2Antlr(fullInput, &dummy, &parserOutputBuffer)
		parserOutput := parserOutputBuffer.String()
		fullExpectedOutput := "parser grammar outputParser;\noptions { tokenVocab=outputLexer; }\n\n " + test.output
		if parserOutput != fullExpectedOutput || err != nil {
			t.Errorf("Peg2Antlr(%s) = %q, %v, want %q, <nil>", test.input, parserOutput, err, test.output)
		}
	}
}
