package engine

import "fmt"

type tokenType uint8

const (
	group tokenType = iota
	bracket
	or
	repeat
	literal
	groupUncaptured
)

type token struct {
	tokentype tokenType
	value     interface{}
}

type parseContent struct {
	pos    int
	tokens []token
}

func parseGroup(regex string, ctx *parseContent) {
	ctx.pos++
	for regex[ctx.pos] != ')' {
		process(regex, ctx)
		ctx.pos += 1
	}
}

func parseBracket(regex string, ctx parseContent) {
	ctx.pos++
	var literals []string
	for regex[ctx.pos] != ']' {
		ch := regex[ctx.pos]

		if ch == '-' {
			next := regex[ctx.pos+1]
			prev := literals[len(literals)-1][0]
			literals[len(literals)-1] = fmt.Sprintf("%c%c", prev, next)
			ctx.pos++
		} else {
			literals = append(literals, fmt.Sprintf("%c", ch))
		}
		ctx.pos++
	}
	literalsSet := map[uint8]bool{}

	for _, l := range literals {
		for i := l[0]; i <= l[len(l)-1]; i++ {
			literalsSet[i] = true
		}
	}
	ctx.tokens = append(ctx.tokens, token{
		tokentype: bracket,
		value:     literals,
	})
}

func parseOr(regex string, ctx *parseContent) {

}

func process(regex string, ctx *parseContent) {
	ch := regex[ctx.pos]
	switch ch {
	case '(':
		groupCtx := &parseContent{
			pos:    ctx.pos,
			tokens: []token{},
		}
		parseGroup(regex, groupCtx)
		ctx.tokens = append(ctx.tokens, token{
			tokentype: group,
			value:     groupCtx.tokens,
		})
		break
	case '[':
		parseBracket(regex, ctx)
		break
	case '|':
		parseOr(regex, ctx)
		break
	case '*', '+', '?':
		parseRepeat(regex, ctx)
		break
	case '{':
		parseRepeatSpecified(regex, ctx)
		break
	default:
		t := token{
			tokentype: literal,
			value:     ch,
		}
		ctx.tokens = append(ctx.tokens, t)
	}
}

func parse(regex string) *parseContent {
	ctx := &parseContent{
		pos:    0,
		tokens: []token{},
	}
	for ctx.pos < len(regex) {
		process(regex, ctx)
		ctx.pos++
	}
	return ctx
}
