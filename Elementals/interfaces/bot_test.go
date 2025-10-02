package main

import "testing"

func TestEnglishBotGreeting(t *testing.T) {
	bot := englishBot{}
	expected := "Hi There!"
	actual := bot.getGreeting()

	if actual != expected {
		t.Errorf("Expected greeting %s, got %s", expected, actual)
	}
}

func TestSpanishBotGreeting(t *testing.T) {
	bot := spanishBot{}
	expected := "¡Hola!"
	actual := bot.getGreeting()

	if actual != expected {
		t.Errorf("Expected greeting %s, got %s", expected, actual)
	}
}

func TestBotInterface(t *testing.T) {
	// Test that both bots implement the interface
	var bot1 bot = englishBot{}
	var bot2 bot = spanishBot{}

	// Test english bot through interface
	greeting1 := bot1.getGreeting()
	if greeting1 != "Hi There!" {
		t.Errorf("Expected English greeting, got %s", greeting1)
	}

	// Test spanish bot through interface
	greeting2 := bot2.getGreeting()
	if greeting2 != "¡Hola!" {
		t.Errorf("Expected Spanish greeting, got %s", greeting2)
	}
}
