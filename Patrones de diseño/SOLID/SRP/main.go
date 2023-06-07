package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// SINGLE RESPONSABILITY PRINCIPLE

var (
	entryCount    = 0
	lineSeparator = "\n"
)

type Journal struct {
	entries []string
}

func (j *Journal) String() string {
	return strings.Join(j.entries, lineSeparator)
}

func (j *Journal) AddEntry(text string) int {
	entryCount++

	entry := fmt.Sprintf("%d: %s", entryCount, text)

	j.entries = append(j.entries, entry)

	return entryCount
}

func (j *Journal) RemoveEntry(index int) {
	// ...
}

// separation of concerns

type Persistence struct{}

func (p *Persistence) SaveToFile(j *Journal, filename string) {
	_ = os.WriteFile(filename, []byte(j.String()), fs.ModeAppend)
}

func main() {
	j := Journal{}
	j.AddEntry("I cried today")
	j.AddEntry("I ate a bug")
	fmt.Println(j.String())

	//
	p := Persistence{}

	p.SaveToFile(&j, "journal.txt")
}
