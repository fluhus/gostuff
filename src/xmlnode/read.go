// Provides a hierarchical representation of XML structures.
package xmlnode

import (
	"encoding/xml"
	"io"
)

// Reads all XML data from the given reader and stores it in a head node.
func ReadAll(r io.Reader) (Node, error) {
	// Create root node.
	// Starting with Tag instead of Root, to eliminate type checks when refering
	// to parent nodes during reading. Will be replaced with a Root node at the
	// end.
	result := &Tag{
		nil,
		"",
		nil,
		nil,
	}
	dec := xml.NewDecoder(r)

	// Parse tokens.
	var t xml.Token
	var err error
	current := result
	for t, err = dec.Token(); err == nil; t, err = dec.Token() {
		switch t := t.(type) {
		case xml.StartElement:
			// Create an attribute map.
			attrs := make([]*xml.Attr, len(t.Attr))
			for i, attr := range t.Attr {
				attrs[i] = &xml.Attr{attr.Name, attr.Value}
			}

			// Create child node.
			child := &Tag{
				current,
				t.Name.Local,
				attrs,
				nil,
			}

			current.children = append(current.children, child)
			current = child

		case xml.EndElement:
			current = current.Parent().(*Tag)

		case xml.CharData:
			child := &Text{
				current,
				string(t),
			}

			current.children = append(current.children, child)

		case xml.Comment:
			child := &Comment{
				current,
				string(t),
			}

			current.children = append(current.children, child)

		case xml.ProcInst:
			child := &ProcInst{
				current,
				string(t.Target),
				string(t.Inst),
			}

			current.children = append(current.children, child)

		case xml.Directive:
			child := &Directive{
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

	return &Root{result.children}, nil
}
