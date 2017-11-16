package token

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	ASSIGN    = "="
	PLUS      = "+"
	MINUS			= "-"
	BANG			= "!"
	ASTERISK	= "*"
	SLASH			= "/"
	LT				= "<"
	GT				= ">"
	LBRACE    = "{"
	RBRACE    = "}"
	LPAREN		= "("
	RPAREN		= ")"
	COMMA			= ","
	SEMICOLON	= ";"
	EQ				= "=="
	NOT_EQ		= "!="
	FUNCTION  = "FUNCTION"
	LET       = "LET"
	TRUE			= "TRUE"
	FALSE 		= "FALSE"
	IF 				= "IF"
	ELSE			= "ELSE"
	RETURN 		= "RETURN"
)

// TODO: add filename and line number / column info

var keywords = map[string]TokenType{
	"fn":  		FUNCTION,
	"let": 		LET,
	"true": 	TRUE,
	"false": 	FALSE,
	"if":			IF,
	"else":		ELSE,
	"return":	RETURN,
}

type TokenType string // TODO: might not need to use string, just byte enums
type Token struct {
	Type    TokenType
	Literal string
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
