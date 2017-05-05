package goelement_test

import "github.com/Sjc1000/goelement"

func ExampleNode_FindTagReverse() {
	parent := child.FindTagReverse("a")
	fmt.Println(parent.Data)
}

func ExampleNode_FindTag() {
	child := parent.FindTag("a")
	fmt.Println(parent.Data)
}

func ExampleNode_FindPath() {
	path := goelement.NodePath{Path: "div/h1"}
	element := root.FindPath(path)
	fmt.Println(element.Data)
}

func ExampleNode_FindPathAll() {
	path := goelement.NodePath{Path: "div/a"}
	elements := root.FindPathAll(path)
	for _, element := range elements {
		fmt.Println(element.Data)
	}
}

func ExampleNode_PrintStructure() {
	root := goelement.ParseFromString(html_data)
	root.PrintStructure(0, "  ")
}
