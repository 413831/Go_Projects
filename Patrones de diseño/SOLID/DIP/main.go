package main

import "fmt"

// Dependency Inversion Principle
// HLM should not depend on LLM
// High Level Module (HLM) - Low Level Module (LLM)
// Both should depend on abstractions

type Relationship int

const (
	Parent Relationship = iota
	Child
	Sibling
)

type Person struct {
	name string
}

type Info struct {
	from         *Person
	relationship Relationship
	to           *Person
}

// LLM
type RelationShipBrowser interface {
	FindAllChildrenOf(name string) []*Person
}

type Relationships struct {
	relations []Info
}

func (r *Relationships) FindAllChildrenOf(name string) []*Person {
	result := make([]*Person, 0)

	for i, v := range r.relations {
		if v.relationship == Parent && v.from.name == name {
			result = append(result, r.relations[i].to)
		}
	}

	return result
}

func (r *Relationships) AddParentAndChild(parent, child *Person) {
	r.relations = append(r.relations, Info{
		from:         parent,
		relationship: Parent,
		to:           child,
	})
}

// HLM
type Research struct {
	browser RelationShipBrowser

	// break DIP
	//relationships Relationships
}

func (r *Research) Investigate(criteria func(Info) bool) {
	for _, p := range r.browser.FindAllChildrenOf("John") {
		fmt.Println("John has a child called", p.name)
	}
}

func criteriaJohnParent(rel Info) bool {
	return rel.from.name == "John" &&
		rel.relationship == Parent
}

func main() {
	parent := Person{"John"}
	child1 := Person{"Chris"}
	child2 := Person{"Matt"}

	relationships := Relationships{}
	relationships.AddParentAndChild(&parent, &child1)
	relationships.AddParentAndChild(&parent, &child2)

	r := Research{&relationships}
	r.Investigate(criteriaJohnParent)
}
