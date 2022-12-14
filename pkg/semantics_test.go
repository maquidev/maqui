package maqui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParserMocker struct {
	buf []Expr
	pos int
}

func NewParserMocker(exprs []Expr) *ParserMocker {
	return &ParserMocker{
		buf: exprs,
		pos: 0,
	}
}

func (b *ParserMocker) Do() {
	return
}

func (b *ParserMocker) Get() Expr {
	if len(b.buf) <= b.pos {
		return &EOS{}
	}

	expr := b.buf[b.pos]
	b.pos++

	return expr
}

func (b *ParserMocker) GetFilename() string {
	return "testing"
}

func TestContextAnalyzer(t *testing.T) {
	cases := []struct {
		name   string
		data   []Expr
		expect *AST
	}{
		{
			"VarSumInt",
			[]Expr{
				&FuncDecl{
					Name: "main",
					Body: []Expr{
						&VariableDecl{
							Name: "x",
							Value: &BinaryExpr{
								Operation: BinaryAddition,
								Op1: &LiteralExpr{
									Typ:   LiteralNumber,
									Value: "1",
								},
								Op2: &LiteralExpr{
									Typ:   LiteralNumber,
									Value: "1",
								},
							},
						},
					},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &FuncDecl{
							Name: "main",
							Body: []Expr{
								&VariableDecl{
									Name: "x",
									Value: &BinaryExpr{
										Operation: BinaryAddition,
										Op1: &LiteralExpr{
											Typ:   LiteralNumber,
											Value: "1",
										},
										Op2: &LiteralExpr{
											Typ:   LiteralNumber,
											Value: "1",
										},
									},
									ResolvedType: &BasicType{
										Typ: "int",
									},
								},
							},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{
								"x":    &BasicType{"int"},
								"main": &FuncType{nil, nil},
							},
						},
					},
				},
				Global: &SymbolTable{
					Entries: map[string]Type{
						"main": &FuncType{nil, nil},
					},
				},
			},
		},
		{
			"IncompatibleTypeIntString",
			[]Expr{
				&VariableDecl{
					Name: "x",
					Value: &BinaryExpr{
						Operation: BinaryAddition,
						Op1: &LiteralExpr{
							Typ:   LiteralNumber,
							Value: "1",
						},
						Op2: &LiteralExpr{
							Typ:   LiteralString,
							Value: "text",
						},
					},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &VariableDecl{
							Name: "x",
							Value: &BinaryExpr{
								Operation: BinaryAddition,
								Op1: &LiteralExpr{
									Typ:   LiteralNumber,
									Value: "1",
								},
								Op2: &LiteralExpr{
									Typ:   LiteralString,
									Value: "text",
								},
							},
							ResolvedType: &TypeErr{TypeErrIncompatible},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{
								"x": &TypeErr{TypeErrIncompatible},
							},
							Errors: []CompileError{
								&IncompatibleTypesError{
									Type1: &BasicType{Typ: "int"},
									Type2: &BasicType{Typ: "string"},
								},
							},
						},
					},
				},
				Errors: []CompileError{
					&IncompatibleTypesError{
						Type1: &BasicType{Typ: "int"},
						Type2: &BasicType{Typ: "string"},
					},
				},
				Global: &SymbolTable{
					Entries: map[string]Type{
						"x": &TypeErr{TypeErrIncompatible},
					},
				},
			},
		},
		{
			"FunctionDeclaration",
			[]Expr{
				&FuncDecl{
					Name: "foo",
					Body: []Expr{},
				},
				&FuncCall{
					Name: "foo",
					Args: []Expr{},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &FuncDecl{
							Name: "foo",
							Body: []Expr{},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{
								"foo": &FuncType{nil, nil},
							},
						},
					},
					{
						Expr: &FuncCall{
							Name: "foo",
							Args: []Expr{},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{
								"foo": &FuncType{nil, nil},
							},
						},
					},
				},
				Errors: nil,
				Global: &SymbolTable{
					Entries: map[string]Type{
						"foo": &FuncType{nil, nil},
					},
				},
			},
		},
		{
			"FunctionCallUndefined",
			[]Expr{
				&FuncCall{
					Name: "foo",
					Args: []Expr{},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &FuncCall{
							Name: "foo",
							Args: []Expr{},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{},
							Errors: []CompileError{
								&UndefinedError{
									Name: "foo",
								},
							},
						},
					},
				},
				Errors: []CompileError{
					&UndefinedError{
						Name: "foo",
					},
				},
				Global: NewGlobalSymbolTable(),
			},
		},
		{
			"IdentifierUndefined",
			[]Expr{
				&Identifier{
					Name: "x",
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &Identifier{
							Name: "x",
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{},
							Errors: []CompileError{
								&UndefinedError{
									Name: "x",
								},
							},
						},
					},
				},
				Errors: []CompileError{
					&UndefinedError{
						Name: "x",
					},
				},
				Global: NewGlobalSymbolTable(),
			},
		},
		{
			"UnaryNegativeImmediate",
			[]Expr{
				&UnaryExpr{
					Operation: UnaryNegative,
					Operand: &LiteralExpr{
						Typ:   LiteralNumber,
						Value: "1",
					},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &UnaryExpr{
							Operation: UnaryNegative,
							Operand: &LiteralExpr{
								Typ:   LiteralNumber,
								Value: "1",
							},
						},
						Stab: NewSymbolTable(),
					},
				},
				Errors: nil,
				Global: NewGlobalSymbolTable(),
			},
		},
		{
			"UnaryNegativeString",
			[]Expr{
				&UnaryExpr{
					Operation: UnaryNegative,
					Operand: &LiteralExpr{
						Typ:   LiteralString,
						Value: "foo",
					},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &UnaryExpr{
							Operation: UnaryNegative,
							Operand: &LiteralExpr{
								Typ:   LiteralString,
								Value: "foo",
							},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{},
							Errors: []CompileError{
								&UndefinedUnitaryError{
									Type: &BasicType{"string"},
									Op:   UnaryNegative,
								},
							},
						},
					},
				},
				Errors: []CompileError{
					&UndefinedUnitaryError{
						Type: &BasicType{"string"},
						Op:   UnaryNegative,
					},
				},
				Global: NewSymbolTable(),
			},
		},
		{
			"StringSubtraction",
			[]Expr{
				&BinaryExpr{
					Operation: BinarySubtraction,
					Op1: &LiteralExpr{
						Typ:   LiteralString,
						Value: "foo",
					},
					Op2: &LiteralExpr{
						Typ:   LiteralString,
						Value: "bar",
					},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &BinaryExpr{
							Operation: BinarySubtraction,
							Op1: &LiteralExpr{
								Typ:   LiteralString,
								Value: "foo",
							},
							Op2: &LiteralExpr{
								Typ:   LiteralString,
								Value: "bar",
							},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{},
							Errors: []CompileError{
								&UndefinedOperationError{
									Type: &BasicType{"string"},
									Op:   BinarySubtraction,
								},
							},
						},
					},
				},
				Errors: []CompileError{
					&UndefinedOperationError{
						Type: &BasicType{"string"},
						Op:   BinarySubtraction,
					},
				},
				Global: NewGlobalSymbolTable(),
			},
		},
		{
			"VarIntLiteralSum",
			[]Expr{
				&VariableDecl{
					Name: "x",
					Value: &BinaryExpr{
						Operation: BinaryAddition,
						Op1: &LiteralExpr{
							Typ:   LiteralNumber,
							Value: "1",
						},
						Op2: &LiteralExpr{
							Typ:   LiteralNumber,
							Value: "1",
						},
					},
				},
				&VariableDecl{
					Name: "y",
					Value: &BinaryExpr{
						Operation: BinaryAddition,
						Op1: &LiteralExpr{
							Typ:   LiteralNumber,
							Value: "1",
						},
						Op2: &Identifier{
							Name: "x",
						},
					},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &VariableDecl{
							Name: "x",
							Value: &BinaryExpr{
								Operation: BinaryAddition,
								Op1: &LiteralExpr{
									Typ:   LiteralNumber,
									Value: "1",
								},
								Op2: &LiteralExpr{
									Typ:   LiteralNumber,
									Value: "1",
								},
							},
							ResolvedType: &BasicType{"int"},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{
								"x": &BasicType{"int"},
								"y": &BasicType{"int"},
							},
						},
					},
					{
						Expr: &VariableDecl{
							Name: "y",
							Value: &BinaryExpr{
								Operation: BinaryAddition,
								Op1: &LiteralExpr{
									Typ:   LiteralNumber,
									Value: "1",
								},
								Op2: &Identifier{
									Name: "x",
								},
							},
							ResolvedType: &BasicType{"int"},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{
								"x": &BasicType{"int"},
								"y": &BasicType{"int"},
							},
						},
					},
				},
				Errors: nil,
				Global: &SymbolTable{
					Entries: map[string]Type{
						"x": &BasicType{"int"},
						"y": &BasicType{"int"},
					},
				},
			},
		},
		{
			"VarDeclarationUndefinedOperand",
			[]Expr{
				&VariableDecl{
					Name: "y",
					Value: &BinaryExpr{
						Operation: BinaryAddition,
						Op1: &LiteralExpr{
							Typ:   LiteralNumber,
							Value: "1",
						},
						Op2: &Identifier{
							Name: "x",
						},
					},
				},
			},
			&AST{
				Statements: []*AnnotatedExpr{
					{
						Expr: &VariableDecl{
							Name: "y",
							Value: &BinaryExpr{
								Operation: BinaryAddition,
								Op1: &LiteralExpr{
									Typ:   LiteralNumber,
									Value: "1",
								},
								Op2: &Identifier{
									Name: "x",
								},
							},
							ResolvedType: &TypeErr{TypeErrUndefined},
						},
						Stab: &SymbolTable{
							Entries: map[string]Type{
								"y": &TypeErr{TypeErrUndefined},
							},
							Errors: []CompileError{
								&UndefinedError{
									Name: "x",
								},
							},
						},
					},
				},
				Errors: []CompileError{
					&UndefinedError{
						Name: "x",
					},
				},
				Global: &SymbolTable{
					Entries: map[string]Type{
						"y": &TypeErr{TypeErrUndefined},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			parser := NewParserMocker(c.data)
			analyzer := NewContextAnalyser(parser)

			c.expect.Filename = parser.GetFilename()

			global := NewGlobalSymbolTable()
			analyzer.DefineInto(global)
			global.Errors = nil

			// Bring global definitions into the expected global table
			c.expect.Global.Import(*NewGlobalSymbolTable())
			for _, ae := range c.expect.Statements {
				ae.Stab.Import(*NewGlobalSymbolTable())
			}

			got := analyzer.Do(global)
			if !assert.Equal(t, c.expect, got) {
				assert.Fail(t, "Unexpected")
			}
		})
	}
}

func TestTypeEquals(t *testing.T) {
	tInt1 := &BasicType{"int"}
	tInt2 := &BasicType{"int"}
	tStr := &BasicType{"string"}

	tFunc1 := &FuncType{
		Args: []*ArgumentType{
			{
				Name: "arg1",
				Type: tInt1,
			},
		},
		Returns: []*BasicType{tStr},
	}

	tFunc2 := &FuncType{
		Args: []*ArgumentType{
			{
				Name: "arg1",
				Type: tInt1,
			},
		},
		Returns: []*BasicType{tStr},
	}

	tFunc3 := &FuncType{
		Args: []*ArgumentType{
			{
				Name: "arg1",
				Type: tInt1,
			},
		},
		Returns: []*BasicType{tInt1},
	}

	assert.True(t, tInt1.Equals(tInt2))
	assert.True(t, tInt2.Equals(tInt1))
	assert.False(t, tStr.Equals(tInt1))
	assert.False(t, tInt1.Equals(tStr))
	assert.False(t, tFunc1.Equals(tStr))
	assert.True(t, tFunc1.Equals(tFunc2))
	assert.True(t, tFunc2.Equals(tFunc1))
	assert.False(t, tFunc2.Equals(tFunc3))
	assert.False(t, tFunc1.Equals(tFunc3))
}

func TestTypeString(t *testing.T) {
	tInt := &BasicType{"int"}
	tFunc := &FuncType{
		Args: []*ArgumentType{
			{
				Name: "arg1",
				Type: &BasicType{"string"},
			},
			{
				Name: "arg2",
				Type: &BasicType{"int"},
			},
		},
		Returns: []*BasicType{
			{"string"},
			{"int"},
		},
	}

	assert.Equal(t, "int", tInt.String())
	assert.Equal(t, "func(string, int) string, int", tFunc.String())
}

func TestStabCopy(t *testing.T) {
	stab := &SymbolTable{
		Entries: map[string]Type{
			"foo": &BasicType{"int"},
			"bar": &FuncType{
				Args: []*ArgumentType{
					{
						Name: "arg1",
						Type: &BasicType{"string"},
					},
				},
				Returns: []*BasicType{
					{"string"},
					{"int"},
				},
			},
		},
		Errors: []CompileError{
			UndefinedOperationError{
				Type: &BasicType{"string"},
				Op:   BinarySubtraction,
			},
		},
	}

	assert.Equal(t, stab, stab.Copy())
}
