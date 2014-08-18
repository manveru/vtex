package vtex

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"strconv"
	"strings"
)

type Reader interface {
	ReadRune() (rune, int, error)
	ReadString(byte) (string, error)
	UnreadRune() error
}

type parser struct {
	buf Reader
}

type Element struct {
	Key   string
	Value map[string]interface{}
}

type Vector3 struct {
	X, Y, Z int64
}

type Vector4 struct {
	W, X, Y, Z int64
}

func ParseReader(b Reader) Element {
	p := parser{buf: b}
	return p.parse()
}

func ParseBytes(b []byte) Element {
	p := parser{buf: bytes.NewBuffer(b)}
	return p.parse()
}

func ParseFile(path string) Element {
	fd, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	p := parser{buf: bufio.NewReader(fd)}
	return p.parse()
}

func (p *parser) parse() (value Element) {
	p.onInstruction()
	return p.onElement()
}

func (p *parser) onInstruction() {
	for {
		r, _, err := p.buf.ReadRune()
		if err != nil {
			return
		}
		switch r {
		case '>':
			return
		}
	}
}
func (p *parser) onName() string {
	s := p.onString()
	return s
}
func (p *parser) onMap() map[string]interface{} {
	m := map[string]interface{}{}

	if !p.expect('{') {
		p.fail("expected {")
	}
	var value interface{}
	for {
		if p.expect('}') {
			break
		}
		key := p.onString()
		t := p.onString()

		switch t {
		case "string":
			value = p.onString()
		case "element_array":
			value = p.onElementArray()
		case "vector4":
			value = p.onVector4()
		case "vector3":
			value = p.onVector3()
		case "int":
			value = p.onInt()
		case "bool":
			value = p.onBool()
		case "string_array":
			value = p.onStringArray()
		case "CDmeImageProcessor":
			key, value = t, p.onMap()
		default:
			p.fail(fmt.Errorf("failed to handle: %s", t))
		}

		m[key] = value
	}

	return m
}

func (p *parser) onBool() bool {
	if p.expect('"') {
		if line, err := p.buf.ReadString('"'); err == nil {
			return line[0:len(line)-1] == "1"
		}
	}

	p.fail("expected int")
	return false
}

func (p *parser) onInt() int64 {
	if p.expect('"') {
		if line, err := p.buf.ReadString('"'); err == nil {
			return atoi64(line[0 : len(line)-1])
		}
	}

	p.fail("expected int")
	return 0
}

func (p *parser) onVector3() (v Vector3) {
	if p.expect('"') {
		if line, err := p.buf.ReadString('"'); err == nil {
			parts := strings.SplitN(line[0:len(line)-1], " ", 3)
			v = Vector3{X: atoi64(parts[0]), Y: atoi64(parts[1]), Z: atoi64(parts[2])}
			return v
		}
	}

	p.fail("expected Vector4")
	return v
}

func (p *parser) onVector4() (v Vector4) {
	if p.expect('"') {
		if line, err := p.buf.ReadString('"'); err == nil {
			parts := strings.SplitN(line[0:len(line)-1], " ", 4)
			v = Vector4{W: atoi64(parts[0]), X: atoi64(parts[1]), Y: atoi64(parts[2]), Z: atoi64(parts[3])}
			return v
		}
	}
	p.fail("expected Vector4")
	return v
}

func (p *parser) onString() string {
	line, err := p.buf.ReadString('"')
	if err != nil {
		p.fail(err)
	}
	line, err = p.buf.ReadString('"')
	if err != nil {
		p.fail(err)
	}
	return line[0 : len(line)-1]
}

func (p *parser) onElement() Element {
	key := p.onString()
	value := p.onMap()
	return Element{Key: key, Value: value}
}

func (p *parser) onElementArray() []interface{} {
	a := []interface{}{}
	p.expect('[')
	for {
		if p.expect(']') {
			break
		}
		a = append(a, p.onElement())
	}
	return a
}

func (p *parser) onStringArray() []string {
	a := []string{}
	p.expect('[')
	for {
		if p.expect(']') {
			break
		}
		a = append(a, p.onString())
	}
	return a
}

func (p *parser) onSpace() string {
	space := ""
	for {
		r, _, _ := p.buf.ReadRune()
		switch r {
		case ' ', '\r', '\t', '\n':
			space += string(r)
		default:
			p.buf.UnreadRune()
			return space
		}
	}
}

func (p *parser) expect(expectRune rune) bool {
	p.onSpace()
	r, _, err := p.buf.ReadRune()
	if err != nil {
		p.fail(err)
	}
	if r == expectRune {
		return true
	}
	p.buf.UnreadRune()
	return false
}

func atoi64(a string) int64 {
	i, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (p *parser) fail(err interface{}) {
	fmt.Printf("%#v\n", p.buf)
	panic(err)
}
