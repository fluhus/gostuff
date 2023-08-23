// Package xmlnode provides a hierarchical node representation of XML documents.
// This package wraps encoding/xml and can be used instead of it.
//
// Each node has an underlying concrete type, but calling all functions is
// legal. For example, here is how you can traverse the node tree:
//
//	func traverse(n Node) {
//	  // Text() returns an empty string for non-text nodes.
//	  doSomeTextSearch(n.Text())
//
//	  // Children() returns nil for non-parent nodes.
//	  for _, child := range n.Children() {
//	    traverse(child)
//	  }
//	}
package xmlnode

import (
	"encoding/xml"
	"io"
)

// ReadAll reads all XML data from the given reader and stores it in a root node.
func ReadAll(r io.Reader) (Node, error) {
	// Create root node.
	// Starting with Tag instead of Root, to eliminate type checks when referring
	// to parent nodes during reading. Will be replaced with a Root node at the
	// end.
	result := &tag{
		nil,
		"",
		nil,
		nil,
	}
	dec := xml.NewDecoder(r)

	var t xml.Token
	var err error
	current := result

	// Parse tokens.
	for t, err = dec.Token(); err == nil; t, err = dec.Token() {
		switch t := t.(type) {
		case xml.StartElement:
			// Copy attributes.
			attrs := make([]*xml.Attr, len(t.Attr))
			for i, attr := range t.Attr {
				attrs[i] = &xml.Attr{Name: attr.Name, Value: attr.Value}
			}

			// Create child node.
			child := &tag{
				current,
				t.Name.Local,
				attrs,
				nil,
			}

			current.children = append(current.children, child)
			current = child

		case xml.EndElement:
			current = current.Parent().(*tag)

		case xml.CharData:
			child := &text{
				current,
				string(t),
			}

			current.children = append(current.children, child)

		case xml.Comment:
			child := &comment{
				current,
				string(t),
			}

			current.children = append(current.children, child)

		case xml.ProcInst:
			child := &procInst{
				current,
				string(t.Target),
				string(t.Inst),
			}

			current.children = append(current.children, child)

		case xml.Directive:
			child := &directive{
				current,
				string(t),
			}

			current.children = append(current.children, child)
		}
	}

	// EOF is ok.
	if err != io.EOF {
		return nil, err
	}

	return &root{result.children}, nil
}
