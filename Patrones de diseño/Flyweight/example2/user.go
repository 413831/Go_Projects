package main

import (
	"fmt"
	"strings"
)

type User struct {
	Fullname string
}

func NewUser(fullname string) *User {
	return &User{Fullname: fullname}
}

var allNames []string

type User2 struct {
	names []uint8
}

func NewUser2(fullName string) *User2 {
	getOrAdd := func(s string) uint8 {
		for i := range allNames {
			if allNames[i] == s {
				return uint8(i)
			}
		}
		allNames = append(allNames, s)
		return uint8(len(allNames) - 1)
	}

	result := User2{}
	parts := strings.Split(fullName, " ")

	for _, p := range parts {
		result.names = append(result.names, getOrAdd(p))
	}

	return &result
}

func (u *User2) FullName() string {
	var parts []string

	for _, id := range u.names {
		parts = append(parts, allNames[id])
	}

	return strings.Join(parts, " ")
}

func main() {
	john := NewUser("John Doe")
	//jane := NewUser("Jane Doe")
	//alsoJane := NewUser("Jane Smith")

	fmt.Println(john.Fullname)

	john2 := NewUser2("John Doe")
	fmt.Println(john2.FullName())
}
