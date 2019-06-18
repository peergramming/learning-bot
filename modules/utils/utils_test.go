package utils

import "testing"

func TestSpacify(t *testing.T) {
	got := Spacify("LinusTorvalds")
	if got != "Linus Torvalds" {
		t.Errorf("Spacify(\"LinusTorvalds\") = \"%s\"; want \"Linus Torvalds\"", got)
	}
}

func TestSpacify2(t *testing.T) {
	got := Spacify("RichardMatthewStallman")
	if got != "Richard Matthew Stallman" {
		t.Errorf("Spacify(\"RichardMatthewStallman\") = \"%s\"; want \"Richard Matthew Stallman\"", got)
	}
}

func TestSpacify3(t *testing.T) {
	got := Spacify("This Is Spaced")
	if got != "This Is Spaced" {
		t.Errorf("Spacify(\"This Is Spaced\") = \"%s\"; want \"This Is Spaced\"", got)
	}
}

func BenchmarkSpacify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spacify("TodayIsTheDay")
	}
}
