package main

// Interface Segregation Principle
type Document struct {
}

type Machine interface {
	Print(d Document)
	Fax(d Document)
	Scan(d Document)
}

type MultiFunctionPrinter struct {
}

func (m MultiFunctionPrinter) Print(d Document) {
	m.Print(d)
}

func (m MultiFunctionPrinter) Fax(d Document) {
	m.Fax(d)
}

func (m MultiFunctionPrinter) Scan(d Document) {
	m.Scan(d)
}

type OldFashionedPrinter struct {
}

func (o OldFashionedPrinter) Print(d Document) {
	o.Print(d)
}

func (o OldFashionedPrinter) Fax(d Document) {
	panic("operation not supported")
}

func (o OldFashionedPrinter) Scan(d Document) {
	panic("operation not supported")
}

// ISP

type Printer interface {
	Print(d Document)
}

type Scanner interface {
	Scan(d Document)
}

type Fax interface {
	Fax(d Document)
}

type MyPrinter struct {
}

func (m MyPrinter) Print(d Document) {

}

type Photocopier struct {
}

func (p Photocopier) Scan(d Document) {

}

func (p Photocopier) Print(d Document) {

}

type MultiFunctionDevice interface {
	Printer
	Scanner
}

// decorator
type MultiFunctionMachine struct {
	printer Printer
	scanner Scanner
}

func (m MultiFunctionMachine) Print(d Document) {
	m.printer.Print(d)
}

func main() {
	ofp := OldFashionedPrinter{}
	ofp.Scan(Document{})
}
