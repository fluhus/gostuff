package xmlnode

import (
	"encoding/xml"
)

// DESIGN NOTE:
// The idea behind the somewhat clumsy interface (rather than something like
// the empty Token interface) is to eliminate some of the need for type checks.
//
// For example, if you need to traverse the node tree, you don't need to check
// if a node is a Tag node before you call its children. Or if you need to
// search for a specific text, you don't need to check which nodes are actually
// text nodes.
//
// You can, of course, check node types if you still need to.

// Represents a single XML node. Can be one of: Root, Tag, Text, Comment,
// ProcInst or Directive.
//
// String methods return empty strings when called on a non-relevant node. For
// example, calling TagName() on a Text node or vise versa. The Children()
// method returns nil for non Tag or Root nodes. The Attr() method returns nil
// for non Tag nodes.
//
// Parent node is nil only in the root node.
type Node interface {
	Parent() Node
	TagName() string
	Attr() []*xml.Attr
	Children() []Node
	Text() string
	Comment() string
	Target() string
	Inst() string
	Directive() string
}

// Represents the root node of an XML tree. Only has children, no other
// properties.
type Root struct {
	children []Node
}

func (n *Root) Parent() Node      { return nil }
func (n *Root) TagName() string   { return "" }
func (n *Root) Attr() []*xml.Attr { return nil }
func (n *Root) Children() []Node  { return n.children }
func (n *Root) Text() string      { return "" }
func (n *Root) Comment() string   { return "" }
func (n *Root) Target() string    { return "" }
func (n *Root) Inst() string      { return "" }
func (n *Root) Directive() string { return "" }

// Represents a start-end element, along with its children.
type Tag struct {
	parent   Node
	tagName  string
	attr     []*xml.Attr
	children []Node
}

func (n *Tag) Parent() Node      { return n.parent }
func (n *Tag) TagName() string   { return n.tagName }
func (n *Tag) Attr() []*xml.Attr { return n.attr }
func (n *Tag) Children() []Node  { return n.children }
func (n *Tag) Text() string      { return "" }
func (n *Tag) Comment() string   { return "" }
func (n *Tag) Target() string    { return "" }
func (n *Tag) Inst() string      { return "" }
func (n *Tag) Directive() string { return "" }

// Represents raw text data, in which XML escape sequences have been replaced
// by the characters they represent.
type Text struct {
	parent Node
	text   string
}

func (n *Text) Parent() Node      { return n.parent }
func (n *Text) TagName() string   { return "" }
func (n *Text) Attr() []*xml.Attr { return nil }
func (n *Text) Children() []Node  { return nil }
func (n *Text) Text() string      { return n.text }
func (n *Text) Comment() string   { return "" }
func (n *Text) Target() string    { return "" }
func (n *Text) Inst() string      { return "" }
func (n *Text) Directive() string { return "" }

// A Comment represents an XML comment of the form <!--comment-->. The string
// does not include the <!-- and --> comment markers.
type Comment struct {
	parent  Node
	comment string
}

func (n *Comment) Parent() Node      { return n.parent }
func (n *Comment) TagName() string   { return "" }
func (n *Comment) Attr() []*xml.Attr { return nil }
func (n *Comment) Children() []Node  { return nil }
func (n *Comment) Text() string      { return "" }
func (n *Comment) Comment() string   { return n.comment }
func (n *Comment) Target() string    { return "" }
func (n *Comment) Inst() string      { return "" }
func (n *Comment) Directive() string { return "" }

// Represents an XML processing instruction of the form <?target inst?>.
type ProcInst struct {
	parent Node
	target string
	inst   string
}

func (n *ProcInst) Parent() Node      { return n.parent }
func (n *ProcInst) TagName() string   { return "" }
func (n *ProcInst) Attr() []*xml.Attr { return nil }
func (n *ProcInst) Children() []Node  { return nil }
func (n *ProcInst) Text() string      { return "" }
func (n *ProcInst) Comment() string   { return "" }
func (n *ProcInst) Target() string    { return n.target }
func (n *ProcInst) Inst() string      { return n.inst }
func (n *ProcInst) Directive() string { return "" }

// Represents an XML directive of the form <!text>. The string does not include
// the <! and > markers.
type Directive struct {
	parent    Node
	directive string
}

func (n *Directive) Parent() Node      { return n.parent }
func (n *Directive) TagName() string   { return "" }
func (n *Directive) Attr() []*xml.Attr { return nil }
func (n *Directive) Children() []Node  { return nil }
func (n *Directive) Text() string      { return "" }
func (n *Directive) Comment() string   { return "" }
func (n *Directive) Target() string    { return "" }
func (n *Directive) Inst() string      { return "" }
func (n *Directive) Directive() string { return n.directive }
