/*
Package goelement provides structure and easy methods for parsing HTML.

Example:
	package main

	import (
		"fmt"
		glm "goelement"
	)

	func main() {
		var html = `
		<html>
			<body>
				<p>Test</p>
				<div>
					<h1 class="element">Testing</h1>
					<h1>Other</h1>
					<h2>Something</h2>
				</div>
				<h1 class="outer">Wooo!</h1>
				<img class="element"/>
			</body>
		</html>
		`
		root := glm.ParseFromString(html)
		path := &glm.NodePath{Path: "body"}

		fmt.Println(root.FindPath(path).FlattenChildren())
	}
*/
package goelement

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

// NodePath contains information about a Node to find. This is passed to FindPath
// and other functions to let them know what to find.
type NodePath struct {
	Path  string
	ID    string
	Class string
}

// Node is a simple structure containing information about each HTML node.
type Node struct {
	html.Token
	Parent     *Node
	Children   []*Node
	Attributes map[string]*html.Attribute
}

// FindTagReverse finds a tag but goes through parents instead of children.
func (node *Node) FindTagReverse(tag string) *Node {
	if node.Data == tag {
		return node
	}
	if node.Parent == nil {
		return nil
	}
	return node.Parent.FindTagReverse(tag)
}

// FindTag finds and returns a child if it has the tagname specified.
func (node *Node) FindTag(tag string) *Node {
	if node.Data == tag {
		return node
	}
	for _, child := range node.Children {
		result := child.FindTag(tag)
		if result != nil {
			return result
		}
	}
	return nil
}

// FindPath finds a single element with a specific path.
func (node *Node) FindPath(nodePath *NodePath) *Node {
	path := nodePath.Path
	class := nodePath.Class
	ID := nodePath.ID
	if node.MatchesPath(path) == true && node.HasClass(class) && node.HasID(ID) {
		return node
	}
	for _, child := range node.Children {
		result := child.FindPath(nodePath)
		if result != nil {
			return result
		}
	}
	return nil
}

// FindPathAll finds all of the elements that match the given path.
func (node *Node) FindPathAll(nodePath *NodePath) []*Node {
	var nodes []*Node
	node.findAll(nodePath, &nodes)
	return nodes
}

// findAll goes through the tree finding any element that matches the NodePath
func (node *Node) findAll(nodePath *NodePath, nodes *[]*Node) {
	path := nodePath.Path
	class := nodePath.Class
	ID := nodePath.ID
	if node.MatchesPath(path) == true && node.HasClass(class) && node.HasID(ID) {
		*nodes = append(*nodes, node)
	}

	for _, child := range node.Children {
		child.findAll(nodePath, nodes)
	}
}

// HasClass checks if a Node has a class.
func (node *Node) HasClass(class string) bool {
	if class == "" {
		return true
	}
	check, ok := node.Attributes["class"]
	if !ok {
		return false
	}
	return check.Val == class
}

// HasID checks if a Node has an ID.
func (node *Node) HasID(ID string) bool {
	if ID == "" {
		return true
	}
	check, ok := node.Attributes["id"]
	if !ok {
		return false
	}
	return check.Val == ID
}

// Path generates the HTML path of the Node.
func (node *Node) Path() string {
	return node.createPath("")
}

// FlattenChildren stomps all of the children and grandchildren into a single slice.
//
// Do you have too many grandchildren? Can't remember their names and attributes? Look no further! FlattenChildren() has what you're looking for! Flatten any sized family tree in mere moments! Yes you heard correctly! Any size, In just moments! Call 1800 FLATTEN now for a free trial!
func (node *Node) FlattenChildren() []*Node {
	var nodes []*Node
	for _, child := range node.Children {
		child.getChildTree(&nodes)
	}
	return nodes
}

// getChildTree does most of the heavy lifting for FlattenChildren()
func (node *Node) getChildTree(nodes *[]*Node) {
	*nodes = append(*nodes, node)
	for _, child := range node.Children {
		child.getChildTree(nodes)
	}
}

// createPath does the heavy lifting for Path()
func (node *Node) createPath(path string) string {
	path = fmt.Sprintf("/%s%s", node.Data, path)
	if node.Parent == nil {
		return path
	}
	return node.Parent.createPath(path)
}

// PrintStructure prints the structure of the Nodes children.
func (node *Node) PrintStructure(indent int, character string) {
	for i := 0; i < indent; i++ {
		fmt.Print(character)
	}
	fmt.Println(node.Data)
	for _, child := range node.Children {
		child.PrintStructure(indent+1, character)
	}
	if len(node.Children) == 0 {
		return
	}
	for i := 0; i < indent; i++ {
		fmt.Print(character)
	}
	fmt.Println(node.Data)
}

// MatchesPath checks if the current node matches the path string.
func (node *Node) MatchesPath(path string) bool {
	if path == "" {
		return true
	}

	split := strings.Split(path, "/")
	current := split[len(split)-1]
	directChild := strings.HasPrefix(current, ".")
	current = strings.TrimPrefix(current, ".")

	if len(split) == 1 {
		return current == node.Data
	} else if current != node.Data {
		return false
	}

	newPath := strings.Join(split[:len(split)-1], "/")

	if directChild == true {
		return node.Parent.MatchesPath(newPath)
	}

	parentName := strings.TrimPrefix(split[len(split)-2], ".")
	parent := node.FindTagReverse(parentName)

	if parent != nil {
		return parent.MatchesPath(newPath)
	}

	return false
}

// NewNode creates a new Node instance.
func newNode(token html.Token, parent *Node) *Node {
	attrs := make(map[string]*html.Attribute)
	for _, attr := range token.Attr {
		attrs[attr.Key] = &attr
	}
	node := &Node{Token: token, Parent: parent, Attributes: attrs}
	return node
}

// dive through the HTML returning the root node.
func dive(tokenizer *html.Tokenizer) *Node {
	var parent *Node
	var root *Node

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return root
		}
		current := tokenizer.Token()
		switch tokenType {
		case html.StartTagToken:
			node := newNode(current, parent)
			if root == nil {
				root = node
			} else {
				parent.Children = append(parent.Children, node)
			}
			parent = node
		case html.SelfClosingTagToken:
			node := newNode(current, parent)
			parent.Children = append(parent.Children, node)
		case html.EndTagToken:
			tag := parent.FindTagReverse(current.Data)
			if tag != nil {
				parent = tag.Parent
			}
		}
	}
}

// ParseFromURL parses HTML from a website.
func ParseFromURL(url string) (*Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	tokenizer := html.NewTokenizer(response.Body)
	return dive(tokenizer), nil
}

// ParseFromString parses HTML data from a string.
func ParseFromString(data string) *Node {
	reader := strings.NewReader(data)
	tokenizer := html.NewTokenizer(reader)
	return dive(tokenizer)
}
