package prompt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	istrings "github.com/balaji01-4d/go-prompt/strings"
)

var ErrUnsupportedChromaLanguage = errors.New("unsupported chroma language")

func NewChromaLexer(language, style string) (Lexer, error) {
	l := lexers.Get(language)
	if l == nil {
		return nil, fmt.Errorf("%w: %s (check Chroma lexers documentation)", ErrUnsupportedChromaLanguage, language)
	}
	l = chroma.Coalesce(l)
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}
	return &ChromaLexer{
		lexer: l,
		style: s,
	}, nil
}

type ChromaLexer struct {
	lexer  chroma.Lexer
	style  *chroma.Style
	tokens []chroma.Token
	idx    int
	offset istrings.ByteNumber
}

func (l *ChromaLexer) Init(input string) {
	l.idx = 0
	l.offset = 0
	iterator, err := l.lexer.Tokenise(nil, input)
	if err != nil {
		l.tokens = nil
		return
	}
	l.tokens = iterator.Tokens()
}

func (l *ChromaLexer) Next() (Token, bool) {
	if l.idx >= len(l.tokens) {
		return nil, false
	}
	t := l.tokens[l.idx]
	l.idx++

	start := l.offset
	length := istrings.Len(t.Value)
	end := start + length
	l.offset += length

	return &ChromaToken{
		token: t,
		style: l.style,
		first: start,
		last:  end - 1, // Last byte index is inclusive
	}, true
}

type ChromaToken struct {
	token chroma.Token
	style *chroma.Style
	first istrings.ByteNumber
	last  istrings.ByteNumber
}

func (t *ChromaToken) Color() Color                          { return DefaultColor }
func (t *ChromaToken) BackgroundColor() Color                { return DefaultColor }
func (t *ChromaToken) DisplayAttributes() []DisplayAttribute { return nil }
func (t *ChromaToken) FirstByteIndex() istrings.ByteNumber   { return t.first }
func (t *ChromaToken) LastByteIndex() istrings.ByteNumber    { return t.last }

func (t *ChromaToken) ANSI() string {
	entry := t.style.Get(t.token.Type)
	var b strings.Builder
	if entry.Bold == chroma.Yes {
		b.WriteString("\x1b[1m")
	}
	if entry.Underline == chroma.Yes {
		b.WriteString("\x1b[4m")
	}
	if entry.Italic == chroma.Yes {
		b.WriteString("\x1b[3m")
	}
	if entry.Colour.IsSet() {
		b.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", entry.Colour.Red(), entry.Colour.Green(), entry.Colour.Blue()))
	}
	if entry.Background.IsSet() {
		b.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm", entry.Background.Red(), entry.Background.Green(), entry.Background.Blue()))
	}
	return b.String()
}
