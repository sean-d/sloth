package parser

import (
	"github.com/sean-d/sloth/ast"
	"github.com/sean-d/sloth/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
   let x = 5;
   let y = 10;
   let foobar = 12345;
   `
	lex := lexer.New(input)
	parse := New(lex)

	program := parse.ParseProgram()
	checkParserErrors(t, parse)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !helperTestLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func helperTestLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
   return 5;
   return 10;
   return 993322;
   `
	lex := lexer.New(input)
	parse := New(lex)
	program := parse.ParseProgram()
	checkParserErrors(t, parse)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	for _, satement := range program.Statements {
		returnStatement, ok := satement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("satement not *ast.ReturnStatement. got=%T", satement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStatement.TokenLiteral())
		}
	}
}

/*
TestIdentifierExpression parses our input foobar;, checks the parser for errors, makes an assertion about the number
of statements in the *ast.Program node and then checks that the only statement in program.Statements is an *ast.ExpressionStatement.
Then we check that the *ast.ExpressionStatement.Expression is an *ast.Identifier.
Finally we check that our identifier has the correct value of "foobar".
*/
func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lex := lexer.New(input)
	parse := New(lex)

	program := parse.ParseProgram()
	checkParserErrors(t, parse)

	if len(program.Statements) != 1 {
		t.Fatalf("program lacks enough statements: got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement. got %d", program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", statement.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

// checkParserErrors checks the parser for errors and if it has any it prints them as test errors and stops the execution of the current test.
func checkParserErrors(t *testing.T, parse *Parser) {
	errors := parse.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}
	t.FailNow()
}
