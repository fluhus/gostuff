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
	case *root:
		result := ""
		for _, child := range n.children {
			result += nodeToString(child)
		}
		return result

	case *tag:
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

	case *text:
		return n.text

	case *comment:
		return "<!--" + n.comment + "-->"

	case *procInst:
		return "<?" + n.target + " " + n.inst + "?>"

	case *directive:
		return "<!" + n.directive + ">"

	default:
		panic("Unknown node type.")
	}
}
