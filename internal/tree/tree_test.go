package tree_test

import (
	"bytes"
	"dbdb/internal/storage"
	"dbdb/internal/tree"
	"os"
	"testing"
)

func setupTestStore(t *testing.T) *storage.Physical {
	file, err := os.CreateTemp("", "dbdb_tree_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	file.Close()
	t.Cleanup(func() { os.Remove(file.Name()) })

	store, err := storage.OpenPhysical(file.Name())
	if err != nil {
		t.Fatalf("Failed to open physical store: %v", err)
	}
	return store
}

func TestBTreeInsertAndBranching(t *testing.T) {
	store := setupTestStore(t)

	bTree := tree.NewBTree(store, tree.NodeRef{})

	rootRef, err := bTree.Set("M", []byte("middle"))
	if err != nil {
		t.Fatalf("Failed to set root: %v", err)
	}
	bTree.Root = rootRef // Update active tree pointer

	newRoot, err := bTree.Set("A", []byte("left"))
	if err != nil {
		t.Fatalf("Failed to set left child: %v", err)
	}
	bTree.Root = newRoot

	finalRoot, err := bTree.Set("Z", []byte("right"))
	if err != nil {
		t.Fatalf("Failed to set right child: %v", err)
	}
	bTree.Root = finalRoot

	// Test Retrieval Matrix
	tests := []struct {
		Key      string
		Expected string
	}{
		{"M", "middle"},
		{"A", "left"},
		{"Z", "right"},
	}

	for _, tc := range tests {
		val, err := bTree.Get(tc.Key)
		if err != nil {
			t.Errorf("Failed to get key %s: %v", tc.Key, err)
		}
		if !bytes.Equal(val, []byte(tc.Expected)) {
			t.Errorf("Key %s: Expected %s, got %s", tc.Key, tc.Expected, val)
		}
	}

	_, err = bTree.Get("missing")
	if err == nil {
		t.Errorf("Expected error for missing key, got nil")
	}
}
