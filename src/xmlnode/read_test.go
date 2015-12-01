package xmlnode

import (
	"strings"
	"testing"
)

func TestReadAll(t *testing.T) {
	text := "<hello>world<?proc inst?><goodbye until=\"later\"><!mydirective>world<!--My comment--></goodbye></hello>"
	node, err := ReadAll(strings.NewReader(text))

	if err != nil {
		t.Fatal(err)
	}

	result := nodeToString(node)
	if result != text {
		t.Fatal("expected:", text, "actual:", result)
	}
}

// Inefficient stringifier for testing.
func nodeToString(n Node) string {
	switch n := n.(type) {
	case *Root:
		result := ""
		for _, child := range n.children {
			result += nodeToString(child)
		}
		return result

	case *Tag:
		result := "<" + n.tagName
		for _, attr := range n.attr {
			result += " " + attr.Name.Local + "=\"" + attr.Value + "\""
		}
		result += ">"
		for _, child := range n.children {
			result += nodeToString(child)
		}
		result += "</" + n.tagName + ">"
		return result

	case *Text:
		return n.text

	case *Comment:
		return "<!--" + n.comment + "-->"

	case *ProcInst:
		return "<?" + n.target + " " + n.inst + "?>"

	case *Directive:
		return "<!" + n.directive + ">"

	default:
		panic("Unknown node type.")
	}
}
