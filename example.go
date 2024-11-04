package main

import (
	"github.com/vdombrovski/xmlparser/xmlparser"
	"fmt"
)

func main() {
	root, err := xmlparser.Parse("sample-data.xml")
	if err != nil {
		panic(err)
	}
	printTree(root, "")
}

func printTree(n *xmlparser.Node, indent string) {
	attrs := ""
	for ak, av := range(n.Attrs) {
		attrs += ak + "=" + av + " "
	}
	fmt.Println(indent + n.Tag, attrs)
	
	for _, c := range n.Children {
		printTree(c, indent + "  ")
	}
}