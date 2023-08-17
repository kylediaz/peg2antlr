# peg2antlr

peg and ANTLR are tools that generate parsers from grammar files. This program
translates the syntax from pointlander/peg-style-grammar to antlr g4 style.

## Usage

```bash
go run ./ peg_input.peg [antlr_output]
```

Creates `antlr_outputLexer.g4` and `antlr_outputLexer.g4`

## How it works


## Credits

`peg_grammar.peg`, the grammar file used to parse .peg files, is shamelessly stolen from the [pointlander/peg repo](https://github.com/pointlander/peg/)
