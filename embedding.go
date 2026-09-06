package main

import "fmt"

type Parent struct {
	name string
}

func (p Parent) Greet() {
	fmt.Println("Hello, my name is", p.name)
}

type ChildEmbed struct {
	Parent // anonymous embed -> promoted
	name   string
}

//func (c ChildEmbed) Greet() {
//	fmt.Println("ChildEmbed:", c.name)
//}
//func (c Child) Greet() {
//	fmt.Println("Hello, my name is", c.name, "and my parent's name is", c.p.name)
//}

type Child struct {
	p    Parent
	name string
}

func main7() {
	ce := ChildEmbed{Parent{"parentName"}, "childNameE"}
	//cn := Child{p: Parent{"parentName"}, name: "childNameN"}

	// if Child defines Greet -> both call child's method otherwise call parent method
	ce.Greet()        // ChildEmbed: childNameE
	ce.Parent.Greet() // Parent: parentName
	fmt.Println("-----")

	
}
