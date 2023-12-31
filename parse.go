package jsondsl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/wenooij/bufiog"
)

type parser struct {
	*bufiog.Reader[tokenPos]
}

func (p *parser) consumeToken(t Token) (Pos, error) {
	e, err := p.ReadElem()
	if err != nil {
		return NoPos, err
	}
	if e.Token != t {
		return NoPos, fmt.Errorf("expected token %s (found %s)", t, e.Token)
	}
	return e.Pos, nil
}

func Parse(src string) ([]Node, error) {
	t := &Tokenizer{}
	sc := bufio.NewScanner(strings.NewReader(src))
	sc.Split(t.SplitFunc)
	p := &parser{
		Reader: bufiog.NewReaderSize(&tokenReader{
			t:  t,
			sc: sc,
		}, 64),
	}

	out := []Node{}
	for {
		if _, err := p.Peek(1); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		val, err := p.parseValue()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		out = append(out, val)
	}
	return out, nil
}

func (p *parser) parseValue() (Value, error) {
	es, err := p.Peek(1)
	if err != nil {
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		return nil, err
	}
	switch e := es[0]; e.Token {
	case TokenInvalid:
		return nil, fmt.Errorf("invalid token returned during scan")
	case TokenColon, TokenComma, TokenLParen, TokenRParen, TokenRBrace, TokenRBrack:
		return nil, fmt.Errorf("unexpected token %s at beginning of Value", e.Token)
	case TokenLBrace:
		object, err := p.parseObject()
		if err != nil {
			return nil, err
		}
		return object, nil
	case TokenLBrack:
		array, err := p.parseArray()
		if err != nil {
			return nil, err
		}
		return array, nil
	case TokenNull:
		p.Discard(1)
		return &Null{NullPos: e.Pos}, nil
	case TokenFalse:
		p.Discard(1)
		return &Bool{LitPos: e.Pos}, nil
	case TokenTrue:
		p.Discard(1)
		return &Bool{LitPos: e.Pos, Literal: true}, nil
	case TokenNumber:
		p.Discard(1)
		return &Number{LitPos: e.Pos, Literal: e.Text}, nil
	case TokenIdent:
		v, err := p.parseOperator()
		if err != nil {
			return nil, err
		}
		return v, nil
	case TokenString:
		return p.parseString()
	default:
		return nil, fmt.Errorf("unknown token %s returned during scan", e.Token)
	}
}

func (p *parser) parseArray() (*Array, error) {
	lb, err := p.consumeToken(TokenLBrack)
	if err != nil {
		return nil, fmt.Errorf("%v at start of array", err)
	}
	elems, err := parseList(p, TokenRBrack, p.parseValue)
	if err != nil {
		return nil, fmt.Errorf("%v in array", err)
	}
	rb, err := p.consumeToken(TokenRBrack)
	if err != nil {
		return nil, fmt.Errorf("%v at end of array", err)
	}
	return &Array{LBrack: lb, Elements: elems, RBrack: rb}, nil
}

func (p *parser) parseObject() (*Object, error) {
	lb, err := p.consumeToken(TokenLBrace)
	if err != nil {
		return nil, fmt.Errorf("%v at beginning of object", err)
	}
	members, err := parseList(p, TokenRBrace, p.parseMember)
	if err != nil {
		return nil, fmt.Errorf("%v in object", err)
	}
	rb, err := p.consumeToken(TokenRBrace)
	if err != nil {
		return nil, fmt.Errorf("%v at end of object", err)
	}
	return &Object{LBrace: lb, Members: members, RBrace: rb}, nil
}

func (p *parser) parseString() (*String, error) {
	e, err := p.ReadElem()
	if err != nil {
		return nil, err
	}
	if e.Token != TokenString {
		return nil, fmt.Errorf("expected token %s (found %s)", TokenString, e.Token)
	}
	return &String{Quote: e.Pos, QuotedContent: e.Text}, nil
}

func (p *parser) parseIdent() (*Ident, error) {
	e, err := p.ReadElem()
	if err != nil {
		return nil, err
	}
	if e.Token != TokenIdent {
		return nil, fmt.Errorf("expected token %s (found %s)", TokenIdent, e.Token)
	}
	return &Ident{NamePos: e.Pos, Name: e.Text}, nil
}

func (p *parser) parseMember() (*Member, error) {
	key, err := p.parseValue()
	if err != nil {
		return nil, fmt.Errorf("%v at member key", err)
	}
	colon, err := p.consumeToken(TokenColon)
	if err != nil {
		return nil, fmt.Errorf("%v in object member", err)
	}
	value, err := p.parseValue()
	if err != nil {
		return nil, fmt.Errorf("%v at member Value", err)
	}
	return &Member{Key: key, Colon: colon, Value: value}, nil
}

func (p *parser) parseOperator() (Value, error) {
	id, err := p.parseIdent()
	if err != nil {
		return nil, fmt.Errorf("%v at start of operator", err)
	}
	var opArgs []*OperatorArgs
	for {
		es, err := p.Peek(1)
		if err != nil && err != io.EOF {
			return nil, err
		}
		var args *OperatorArgs
		if len(es) == 0 || es[0].Token != TokenLParen {
			break
		}
		if args, err = p.parseOperatorArgs(); err != nil {
			return nil, err
		}
		opArgs = append(opArgs, args)
	}
	return &Operator{Id: id, Args: opArgs}, nil
}

func (p *parser) parseOperatorArgs() (*OperatorArgs, error) {
	lp, err := p.consumeToken(TokenLParen)
	if err != nil {
		return nil, fmt.Errorf("%v at start of operator arguments", err)
	}
	args, err := parseList(p, TokenRParen, p.parseValue)
	if err != nil {
		return nil, fmt.Errorf("%v at operator arguments", err)
	}
	rp, err := p.consumeToken(TokenRParen)
	if err != nil {
		return nil, fmt.Errorf("%v at end of operator", err)
	}
	return &OperatorArgs{
		LParen:    lp,
		ValueList: args,
		RParen:    rp,
	}, nil
}

// parseList parses a generic list of Nodes as seen in the object, array, and operator specs.
// It parses the contents of the list including TokenComma, but does not consume the provided
// delim.
//
// precondition: delim is one of: TokenRBrack, TokenBrace, or TokenRParen.
func parseList[E Node](p *parser, delim Token, parseFn func() (E, error)) ([]ListElem[E], error) {
	out := []ListElem[E](nil)

	for done := false; !done; {
		es, err := p.Peek(1)
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
		if es[0].Token == delim {
			break
		}
		v, err := parseFn()
		if err != nil {
			return nil, err
		}
		es, err = p.Peek(1)
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
		var comma Pos
		switch es[0].Token {
		case TokenComma:
			comma = es[0].Pos
			p.Discard(1)
		case delim:
			done = true
		default:
			return nil, fmt.Errorf("expected token %s (found %s)", TokenComma, es[0].Token)
		}
		out = append(out, ListElem[E]{Value: v, Comma: comma})
	}
	return out, nil
}
