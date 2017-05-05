# goelement

A simple HTML structure creator based on Golangs inbuilt HTML parser.

This creates nodes that have child / parent information. It also provides some useful searching methods on those nodes such as .FindTag() and .FindPath()


## Installing

`go get github.com/Sjc1000/goelement`


## Usage 

`import "github.com/Sjc1000/goelement"`

Here is a basic example


```Go
package main

import (
	"fmt"
	glm "github.com/Sjc1000/goelement"
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
			<h1 class="outer something">Wooo!</h1>
			<img class="element"/>
		</body>
	</html>
	`
	path := &glm.NodePath{Path: "h1", Class: "element"}
	root := glm.ParseFromString(html)
	fmt.Println(root.FindPath(path).Path())
}
```
