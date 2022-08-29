package maqui

type Compiler struct{}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(filename string) (error, []CompileError) {
	lexer, err := NewLexer(filename)
	if err != nil {
		return err, nil
	}

	parser := NewParser(lexer)
	analyzer := NewContextAnalyser(parser)

	global := NewGlobalSymbolTable()
	analyzer.Define(global)
	ast := analyzer.Do(global)

	return nil, ast.Errors
}
