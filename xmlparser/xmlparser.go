package xmlparser

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Flags struct {
	open           bool
	content        bool
	attrKey        bool
	attrValueStart bool
	attrValue      bool
	closing        bool
	startClosing   bool
}

type Node struct {
	Parent   *Node
	Children []*Node
	Tag      string
	Content  string
	Attrs    map[string]string
	f        Flags
	buf      string
	abuf     string
}

func newNode(parent *Node) *Node {
	return &Node{
		Parent:   parent,
		Children: []*Node{},
		Attrs:    map[string]string{},
	}
}

func Parse(path string) (*Node, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := make([]byte, 1)
	var idx uint64
	var curr *Node
	for {
		_, err := f.Read(b)
		idx++
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch string(b) {
		case " ":
			if curr.f.open && !curr.f.attrValue {
				if !curr.f.attrKey {
					if curr.Tag == "" {
						curr.Tag = curr.buf
						curr.buf = ""
					}
					curr.f.attrKey = true
				}
				continue
			}
			curr.buf += string(b)
		case "\n":
			if curr.f.open {
				return nil, fmt.Errorf("Syntax error: Unexpected '\\n' in tag '%s' on char %d", curr.Tag, idx)
			}
			curr.buf += string(b)
		case "=":
			if curr.f.open && curr.f.attrKey {
				curr.f.attrKey = false
				curr.f.attrValueStart = true
				curr.abuf = curr.buf
				curr.buf = ""
			}
		case "\"":
			if curr.f.attrValueStart {
				curr.f.attrValueStart = false
				curr.f.attrValue = true
				continue
			}
			if curr.f.attrValue {
				curr.Attrs[curr.abuf] = curr.buf
				curr.abuf = ""
				curr.buf = ""
				curr.f.attrValue = false
			}
		case "<":
			if curr == nil {
				curr = newNode(nil)
				curr.f.open = true
				continue
			}
			if curr.f.open {
				return nil, fmt.Errorf("Syntax error: Unexpected '<' in tag '%s' on char %d", curr.Tag, idx)
			}
			if curr.f.attrValueStart {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag '%s' after key '%s' on char %d", string(b), curr.Tag, curr.abuf, idx)
			}
			if curr.f.attrKey {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag '%s' in key name '%s' on char %d", string(b), curr.Tag, curr.buf, idx)
			}
			if curr.f.content {
				curr.f.startClosing = true
			}
		case "/":
			if curr.f.attrValueStart {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag '%s' after key '%s' on char %d", string(b), curr.Tag, curr.abuf, idx)
			}
			if curr.f.attrKey {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag '%s' in key name '%s' on char %d", string(b), curr.Tag, curr.buf, idx)
			}
			if curr.f.startClosing {
				curr.Content = curr.buf
				curr.buf = ""
				curr.f.content = false
				curr.f.closing = true
				curr.f.startClosing = false
			}
		case ">":
			if curr.f.attrValueStart {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag '%s' after key '%s' on char %d", string(b), curr.Tag, curr.abuf, idx)
			}
			if curr.f.attrKey {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag '%s' in key name '%s' on char %d", string(b), curr.Tag, curr.buf, idx)
			}
			if curr.f.open {
				if len(curr.Attrs) < 1 {
					curr.Tag = curr.buf
					curr.buf = ""
				}
				curr.f.open = false
				curr.f.content = true
			}
			if curr.f.closing {
				if curr.buf != curr.Tag {
					return nil, fmt.Errorf("Syntax error: Tag '%s' is not properly closed (got '%s') on char %d", curr.Tag, curr.buf, idx)
				}
				if curr.Parent == nil {
					return curr, nil
				}
				curr.Parent.Children = append(curr.Parent.Children, curr)
				curr = curr.Parent
			}
		default:
			if curr.f.attrValueStart {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag '%s' after key '%s' on char %d", string(b), curr.Tag, curr.abuf, idx)
			}
			if curr.f.attrKey && !strings.Contains("abcdefghijklmnopqrstuvwxyz_-0123456789", string(b)) {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag '%s' in key name '%s' on char %d", string(b), curr.Tag, curr.buf, idx)
			}
			if curr.f.open && !curr.f.attrValue && !strings.Contains("abcdefghijklmnopqrstuvwxyz_-0123456789", string(b)) {
				return nil, fmt.Errorf("Syntax error: Unexpected '%s' in tag name '%s' on char %d", string(b), curr.buf, idx)
			}
			if curr.f.startClosing {
				curr.f.startClosing = false
				curr = newNode(curr)
				curr.f.open = true
			}
			curr.f.startClosing = false
			curr.buf += string(b)
		}
	}
	return nil, errors.New("Syntax error")
}
