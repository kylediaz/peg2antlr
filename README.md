# peg2antlr

peg and ANTLR are tools that generate parsers from grammar files. This program
translates the syntax from pointlander/peg-style grammar to antlr g4 style.

## Usage

```bash
go run ./ peg_input.peg [antlr_output]
```

Creates `antlr_outputLexer.g4` and `antlr_outputLexer.g4`

## How it works

1. The [peg](https://github.com/pointlander/peg/) repository has a grammar for the peg language. I used that to generate a parser for .peg files.
2. Parse the peg file, turning it into an AST
3. Traverse the AST and index the tokens used in the peg grammar file so we can know what Lexer definitions we need.
4. Compile the AST into a ANTLR-style grammar file using various rules.
    * For example, grammar rule definitions in peg are `label <- expression` while in ANTLR it's `label: expression`. This can easily be translated by switching out the symbol.
    * This rules can be more advanced, such as how it handles regex. ANTLR's regex engine
    has less features than peg's, so it requires a little more elbow grease.

## Credits

`peg_grammar.peg`, the grammar file used to parse .peg files, is shamelessly stolen from the [pointlander/peg repo](https://github.com/pointlander/peg/)
