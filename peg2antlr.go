package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/pointlander/peg/tree"
)

func main() {
	// Check if the correct number of arguments are provided
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run main.go peg_input.peg [antlr_output.g4]")
		return
	}

	pegFileName := os.Args[1]
	antlrFileName := "output"

	if len(os.Args) == 3 {
		antlrFileName = os.Args[2]
	}

	pegContent, err := os.ReadFile(pegFileName)
	if err != nil {
		fmt.Println("Error reading peg_input.peg:", err)
		return
	}

	var lexerOutput, parserOutput bytes.Buffer

	err = Peg2Antlr(string(pegContent), &lexerOutput, &parserOutput)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = writeBufferToFile(lexerOutput, antlrFileName+"Lexer.g4")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writeBufferToFile(parserOutput, antlrFileName+"Parser.g4")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Success!")
}

func Peg2Antlr(input string, lexerOutput, parserOutput *bytes.Buffer) error {
	p := &Peg{
		Tree:   tree.New(false, false, false),
		Buffer: input,
	}
	p.Init(Pretty(true), Size(1<<15))
	if err := p.Parse(); err != nil {
		return err
	}
	p.Execute()

	ast := p.AST()
	cleanAst(ast)
	ast = remove(ast, ruleComment)

	literals := getLiterals(ast, input)
	labels := labelTokens(literals, input)
	outputLexer(labels, lexerOutput)

	peg2AntlrImpl(ast, input, parserOutput)

	return nil
}

// Removes all top-level statements that aren't comments or definitions
func cleanAst(node *node32) {
	if node == nil {
		return
	}
	var clean func(node *node32) *node32
	clean = func(node *node32) *node32 {
		if node == nil {
			return nil
		}
		switch node.pegRule {
		case ruleComment, ruleDefinition:
			node.next = clean(node.next)
			return node
		default:
			return clean(node.next)
		}
	}
	node.up = clean(node.up)
}

////////////////////////////////////////////////////////////////////////////////
// Lexer
////////////////////////////////////////////////////////////////////////////////

func getLiterals(node *node32, rawInput string) []*node32 {
	isLiteral := func(n *node32) bool { return n.pegRule == ruleLiteral }
	literals := collect(node, isLiteral)
	return literals
}

// Outputs token literal value -> label name
func labelTokens(literals []*node32, rawInput string) map[string]string {
	labels := make(map[string]string, len(literals))
	for _, l := range literals {
		literalValue := getRawNodeValue(l, rawInput)
		if _, ok := labels[literalValue]; !ok {
			tokenLabel := generateTokenLabel(literalValue[1 : len(literalValue)-2])
			labels[literalValue] = tokenLabel
		}
	}
	return labels
}

func outputLexer(labels map[string]string, output *bytes.Buffer) {
	output.WriteString("lexer grammar Lexer;\n\n")
	for literal, label := range labels {
		definition := fmt.Sprintf("%s: %s;\n", label, literal)
		output.WriteString(definition)
	}
}

////////////////////////////////////////////////////////////////////////////////
// Transforming the Definitions to ANTLR
////////////////////////////////////////////////////////////////////////////////

func Peg2AntlrParser(input string) (string, error) {
	p := &Peg{
		Tree:   tree.New(false, false, false),
		Buffer: input,
	}
	p.Init(Pretty(true), Size(1<<15))
	if err := p.Parse(); err != nil {
		return "", err
	}
	p.Execute()

	var output bytes.Buffer
	output.WriteString("parser grammar outputParser;\noptions { tokenVocab=outputLexer; }\n\n")
	ast := p.AST()

	cleanAst(ast)
	peg2AntlrImpl(ast, input, &output)

	return output.String(), nil
}

func peg2AntlrImpl(node *node32, rawInput string, outputBuffer *bytes.Buffer) {
	for node != nil {
		peg2AntlrImpl2(node, rawInput, outputBuffer)
		node = node.next
	}
}

func peg2AntlrImpl2(node *node32, rawInput string, outputBuffer *bytes.Buffer) {
	_print := func(format string, a ...interface{}) { fmt.Fprintf(outputBuffer, format, a...) }
	rawValue := getRawNodeValue(node, rawInput)
	_printRaw := func() { _print(rawValue) }

	_descend := func() { peg2AntlrImpl(node.up, rawInput, outputBuffer) }

	switch node.pegRule {
	// Skip
	case ruleBegin, ruleEnd, ruleAction:
		return
	// Raw prints
	case ruleOpen, ruleClose, ruleLiteral, rulePlus, ruleDot,
		ruleQuestion, ruleClass, ruleAnd:
		_printRaw()
	case ruleDefinition:
		_descend()
		_print(";\n")
	case ruleIdentifier:
		value := getRawNodeValue(node.up, rawInput) // Underlying PegText
		if value == "_" {
			value = "UNDERSCORE"
		}
		_print(" %s", value)
	case ruleExpression:
		var expressionOutput bytes.Buffer
		peg2AntlrImpl(node.up, rawInput, &expressionOutput)
		expression := expressionOutput.String()
		expression = strings.Replace(expression, "\n", "", 0)
		_print(expression)
	case ruleNot:
		_print("~")
	case ruleLeftArrow:
		_print(":")
	case ruleSlash:
		_print(" | ")
	case ruleStar:
		_print("*")
		_descend()
	case ruleComment:
		_print("// %s", rawValue[1:])
	case rulePegText:
		value := getRawNodeValue(node, rawInput)
		if value == "_" {
			_print(" UNDERSCORE ")
		} else {
			_printRaw()
		}
	default:
		_descend()
	}
}
