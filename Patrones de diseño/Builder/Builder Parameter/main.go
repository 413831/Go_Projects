package main

import "strings"

type (
	Email struct {
		from, to, subject, body string
	}

	EmailBuilder struct {
		email Email
	}

	Build func(builder *EmailBuilder)
)

func (b *EmailBuilder) From(from string) *EmailBuilder {
	if validateEmailAddress(from) {
		b.email.from = from
		return b
	}

	panic("email should contain @")
}

func (b *EmailBuilder) To(to string) *EmailBuilder {
	if validateEmailAddress(to) {
		b.email.to = to
		return b
	}

	panic("email should contain @")
}

func (b *EmailBuilder) Subject(subject string) *EmailBuilder {
	b.email.subject = subject

	return b
}

func (b *EmailBuilder) Body(body string) *EmailBuilder {
	b.email.body = body

	return b
}

func validateEmailAddress(to string) bool {
	if strings.Contains(to, "@") {
		return true
	}

	return false
}

func sendMailImpl(email *Email) {

}

func SendEmail(action Build) {
	builder := EmailBuilder{}
	action(&builder)
	sendMailImpl(&builder.email)
}

func main() {
	SendEmail(func(b *EmailBuilder) {
		b.From("foo@bar.com").
			To("bar@baz.com").
			Subject("Meeting").
			Body("Hello, do ou want to meet?")

	})
}
