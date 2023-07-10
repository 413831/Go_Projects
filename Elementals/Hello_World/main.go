package main

import "fmt"

type School struct {
	courses []string
}

func main() {
	s := School{}
	fmt.Println(len(s.courses))
}
