// Code generated by "stringer -type Token -trimprefix Token"; DO NOT EDIT.

package jsondsl

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokenInvalid-0]
	_ = x[TokenColon-1]
	_ = x[TokenComma-2]
	_ = x[TokenDot-3]
	_ = x[TokenLParen-4]
	_ = x[TokenRParen-5]
	_ = x[TokenLBrace-6]
	_ = x[TokenRBrace-7]
	_ = x[TokenLBrack-8]
	_ = x[TokenRBrack-9]
	_ = x[TokenNull-10]
	_ = x[TokenFalse-11]
	_ = x[TokenTrue-12]
	_ = x[TokenNumber-13]
	_ = x[TokenIdent-14]
	_ = x[TokenString-15]
}

const _Token_name = "InvalidColonCommaDotLParenRParenLBraceRBraceLBrackRBrackNullFalseTrueNumberIdentString"

var _Token_index = [...]uint8{0, 7, 12, 17, 20, 26, 32, 38, 44, 50, 56, 60, 65, 69, 75, 80, 86}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
