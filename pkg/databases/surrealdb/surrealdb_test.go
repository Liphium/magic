package surrealdb

import (
	"testing"
)

func TestLoadSaveRoundtrip(t *testing.T) {
	driver := NewDriver("surrealdb/surrealdb:v3").
		NewDatabase("liphium", "main").
		NewDatabase("liphium", "chat")

	data, err := driver.Save()
	if err != nil {
		t.Fatalf("couldn't save driver: %s", err)
	}

	loaded, err := driver.Load(data)
	if err != nil {
		t.Fatalf("couldn't load driver: %s", err)
	}

	restored := loaded.(*SurrealDriver)
	if restored.Image != driver.Image {
		t.Errorf("expected image %s, got %s", driver.Image, restored.Image)
	}
	if len(restored.Databases) != 2 {
		t.Fatalf("expected 2 databases, got %d", len(restored.Databases))
	}
	if restored.Databases[0] != (SurrealDatabase{Namespace: "liphium", Database: "main"}) {
		t.Errorf("expected first database to be liphium/main, got %v", restored.Databases[0])
	}
	if restored.Databases[1] != (SurrealDatabase{Namespace: "liphium", Database: "chat"}) {
		t.Errorf("expected second database to be liphium/chat, got %v", restored.Databases[1])
	}
}
