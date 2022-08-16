package maqui

import (
	"strings"
	"testing"

	"go.maqui.dev/internal/test"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		data   string
		fail   bool
		expect []Token
	}{
		{
			"func main () {}",
			false,
			[]Token{
				{TokenFunc, "func", nil},
				{TokenIdentifier, "main", nil},
				{TokenOpenParentheses, "(", nil},
				{TokenCloseParentheses, ")", nil},
				{TokenOpenCurly, "{", nil},
				{TokenCloseCurly, "}", nil},
			},
		},
		{
			"//this is a comment\n",
			false,
			[]Token{
				{TokenLineComment, "this is a comment", nil},
			},
		},
		{
			"func main () {\n// this is a comment \n}",
			false,
			[]Token{
				{TokenFunc, "func", nil},
				{TokenIdentifier, "main", nil},
				{TokenOpenParentheses, "(", nil},
				{TokenCloseParentheses, ")", nil},
				{TokenOpenCurly, "{", nil},
				{TokenLineComment, " this is a comment ", nil},
				{TokenCloseCurly, "}", nil},
			},
		},
		{
			"únicódeShouldBeVàlid := 1",
			false,
			[]Token{
				{TokenIdentifier, "únicódeShouldBeVàlid", nil},
				{TokenDeclaration, ":=", nil},
				{TokenNumber, "1", nil},
			},
		},
		{
			"varDeclExpr := \"string\"",
			false,
			[]Token{
				{TokenIdentifier, "varDeclExpr", nil},
				{TokenDeclaration, ":=", nil},
				{TokenString, "string", nil},
			},
		},
		{
			"\"\"",
			false,
			[]Token{
				{TokenString, "", nil},
			},
		},
		{
			"\"unclosed string",
			true,
			nil,
		},
		{
			"@",
			true,
			nil,
		},
	}

	for _, c := range cases {
		r := strings.NewReader(c.data)
		l := NewLexerFromReader(r)

		toks, err := l.RunBlocking()
		if c.fail {
			assert.Error(t, err)
		}

		for i := 0; i < len(toks); i++ {
			toks[i].Loc = nil // ignore meta
		}

		assert.Equal(t, c.expect, toks)

	}
}

// Use a package-level variable to avoid compiler optimisation
var benchResult []Token

func benchmarkLexer(size int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		// Setup
		b.StopTimer()
		data := test.GetRandomTokens(size)
		r := strings.NewReader(data)
		l := NewLexerFromReader(r)

		var err error
		b.StartTimer()

		benchResult, err = l.RunBlocking()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLexer100(b *testing.B) {
	benchmarkLexer(100, b)
}

func BenchmarkLexer1000(b *testing.B) {
	benchmarkLexer(1000, b)
}

func BenchmarkLexer10000(b *testing.B) {
	benchmarkLexer(10000, b)
}

func BenchmarkLexer100000(b *testing.B) {
	benchmarkLexer(100000, b)
}

func BenchmarkLexer1000000(b *testing.B) {
	benchmarkLexer(1000000, b)
}
