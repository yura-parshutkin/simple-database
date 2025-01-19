package storage

import (
	"context"
	"testing"
)

func TestInMemoryEngine_Get(t *testing.T) {
	ctx := context.Background()
	engine := NewInMemoryEngine()
	_ = engine.Set(ctx, "key", "value")

	val, err := engine.Get(ctx, "key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "value" {
		t.Errorf("expected value 'value', got %v", val)
	}

	val, err = engine.Get(ctx, "missing_key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected key to not exist")
	}
}

func TestInMemoryEngine_Set(t *testing.T) {
	ctx := context.Background()
	engine := NewInMemoryEngine()

	err := engine.Set(ctx, "key", "value")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	val, err := engine.Get(ctx, "key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "value" {
		t.Errorf("expected value 'value', got %v", val)
	}
}

func TestInMemoryEngine_Delete(t *testing.T) {
	ctx := context.Background()
	engine := NewInMemoryEngine()
	_ = engine.Set(ctx, "key", "value")

	deleted, err := engine.Delete(ctx, "key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !deleted {
		t.Errorf("expected key to be deleted")
	}

	val, err := engine.Get(ctx, "key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected key to not exist")
	}

	deleted, err = engine.Delete(ctx, "missing_key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if deleted {
		t.Errorf("expected key to not exist")
	}
}
