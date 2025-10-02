package main

import "testing"

func TestPersonCreation(t *testing.T) {
	person := person{
		firstName: "John",
		lastName:  "Doe",
		contact: contactInfo{
			email:   "john@example.com",
			zipCode: 12345,
		},
	}

	if person.firstName != "John" {
		t.Errorf("Expected firstName John, got %s", person.firstName)
	}

	if person.lastName != "Doe" {
		t.Errorf("Expected lastName Doe, got %s", person.lastName)
	}

	if person.contact.email != "john@example.com" {
		t.Errorf("Expected email john@example.com, got %s", person.contact.email)
	}

	if person.contact.zipCode != 12345 {
		t.Errorf("Expected zipCode 12345, got %d", person.contact.zipCode)
	}
}

func TestUpdateName(t *testing.T) {
	person := person{
		firstName: "John",
		lastName:  "Doe",
		contact: contactInfo{
			email:   "john@example.com",
			zipCode: 12345,
		},
	}

	// Test updating name
	person.updateName("Jane")

	if person.firstName != "Jane" {
		t.Errorf("Expected firstName Jane after update, got %s", person.firstName)
	}

	// Ensure other fields remain unchanged
	if person.lastName != "Doe" {
		t.Errorf("Expected lastName to remain Doe, got %s", person.lastName)
	}
}

func TestContactInfo(t *testing.T) {
	contact := contactInfo{
		email:   "test@example.com",
		zipCode: 54321,
	}

	if contact.email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", contact.email)
	}

	if contact.zipCode != 54321 {
		t.Errorf("Expected zipCode 54321, got %d", contact.zipCode)
	}
}
