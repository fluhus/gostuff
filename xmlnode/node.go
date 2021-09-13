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
// You can, of course, check node types (returned by Type()) if you still need
// to.

// DESIGN NOTE (2):
// The Type() function is used instead of exported types to avoid cluttering
// the package's godoc.

// A Node represents a single XML node. Can be one of: Root, Tag, Text, Comment,
// ProcInst or Directive.
//
// String methods return empty strings when called on a non-relevant node. For
// example, calling TagName() on a Text node or vise versa. The Children()
// method returns nil for non Tag or Root nodes. The Attr() method returns nil
// for non Tag nodes.
//
// Parent node is nil only in Root.
type Node interface {
	// Parent of the current node. Nil for root node.
	Parent() Node

	// Tag name of tag nodes of the form <tagname>...</tagname>. Empty for
	// other node types.
	TagName() string

	// Attributes of tag nodes. Nil for other node types.
	Attr() []*xml.Attr

	// Child nodes of root and tag nodes. Nil for other node types.
	Children() []Node

	// Text data of text nodes. Empty for other node types.
	Text() string

	// Comment data of comments of the form <!--comment-->. Does not include the
	// <!-- and --> markers. Empty for other node types.
	Comment() string

	// Target of processing instructions of the form <?target inst?>.
	// Empty for other node types.
	Target() string

	// Instruction of processing instructions of the form <?target inst?>.
	// Empty for other node types.
	Inst() string

	// Directive of the form <!directive>. Does not include the <! and >
	// markers. Empty for other node types.
	Directive() string

	// Type of this node. Returns one of: Root, Tag, Text, Comment, ProcInst
	// or Directive.
	Type() int
}

// Represents the root node of an XML tree. Only has children, no other
// properties.
type root struct {
	children []Node
}

func (n *root) Parent() Node      { return nil }
func (n *root) TagName() string   { return "" }
func (n *root) Attr() []*xml.Attr { return nil }
func (n *root) Children() []Node  { return n.children }
func (n *root) Text() string      { return "" }
func (n *root) Comment() string   { return "" }
func (n *root) Target() string    { return "" }
func (n *root) Inst() string      { return "" }
func (n *root) Directive() string { return "" }
func (n *root) Type() int         { return Root }

// Represents a start-end element, along with its children.
type tag struct {
	parent   Node
	tagName  string
	attr     []*xml.Attr
	children []Node
}

func (n *tag) Parent() Node      { return n.parent }
func (n *tag) TagName() string   { return n.tagName }
func (n *tag) Attr() []*xml.Attr { return n.attr }
func (n *tag) Children() []Node  { return n.children }
func (n *tag) Text() string      { return "" }
func (n *tag) Comment() string   { return "" }
func (n *tag) Target() string    { return "" }
func (n *tag) Inst() string      { return "" }
func (n *tag) Directive() string { return "" }
func (n *tag) Type() int         { return Tag }

// Represents raw text data, in which XML escape sequences have been replaced
// by the characters they represent.
type text struct {
	parent Node
	text   string
}

func (n *text) Parent() Node      { return n.parent }
func (n *text) TagName() string   { return "" }
func (n *text) Attr() []*xml.Attr { return nil }
func (n *text) Children() []Node  { return nil }
func (n *text) Text() string      { return n.text }
func (n *text) Comment() string   { return "" }
func (n *text) Target() string    { return "" }
func (n *text) Inst() string      { return "" }
func (n *text) Directive() string { return "" }
func (n *text) Type() int         { return Text }

// A Comment represents an XML comment of the form <!--comment-->. The string
// does not include the <!-- and --> comment markers.
type comment struct {
	parent  Node
	comment string
}

func (n *comment) Parent() Node      { return n.parent }
func (n *comment) TagName() string   { return "" }
func (n *comment) Attr() []*xml.Attr { return nil }
func (n *comment) Children() []Node  { return nil }
func (n *comment) Text() string      { return "" }
func (n *comment) Comment() string   { return n.comment }
func (n *comment) Target() string    { return "" }
func (n *comment) Inst() string      { return "" }
func (n *comment) Directive() string { return "" }
func (n *comment) Type() int         { return Comment }

// Represents an XML processing instruction of the form <?target inst?>.
type procInst struct {
	parent Node
	target string
	inst   string
}

func (n *procInst) Parent() Node      { return n.parent }
func (n *procInst) TagName() string   { return "" }
func (n *procInst) Attr() []*xml.Attr { return nil }
func (n *procInst) Children() []Node  { return nil }
func (n *procInst) Text() string      { return "" }
func (n *procInst) Comment() string   { return "" }
func (n *procInst) Target() string    { return n.target }
func (n *procInst) Inst() string      { return n.inst }
func (n *procInst) Directive() string { return "" }
func (n *procInst) Type() int         { return ProcInst }

// Represents an XML directive of the form <!text>. The string does not include
// the <! and > markers.
type directive struct {
	parent    Node
	directive string
}

func (n *directive) Parent() Node      { return n.parent }
func (n *directive) TagName() string   { return "" }
func (n *directive) Attr() []*xml.Attr { return nil }
func (n *directive) Children() []Node  { return nil }
func (n *directive) Text() string      { return "" }
func (n *directive) Comment() string   { return "" }
func (n *directive) Target() string    { return "" }
func (n *directive) Inst() string      { return "" }
func (n *directive) Directive() string { return n.directive }
func (n *directive) Type() int         { return Directive }

// Possible return values of Type().
const (
	Root = iota
	Tag
	Text
	Comment
	ProcInst
	Directive
)
