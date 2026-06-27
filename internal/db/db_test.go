package db_test

import (
	"bytes"
	"dbdb/internal/db"
	"os"
	"testing"
)

func TestStorageEngineEndToEnd(t *testing.T) {
	tmpFile := "test_dbdb.db"
	defer os.Remove(tmpFile)

	database, err := db.Open(tmpFile)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	err = database.Set("alpha", []byte("value_1"))
	if err != nil {
		t.Fatalf("Failed to write alpha: %v", err)
	}

	err = database.Set("beta", []byte("value_2"))
	if err != nil {
		t.Fatalf("Failed to write beta: %v", err)
	}

	val, err := database.Get("alpha")
	if err != nil || !bytes.Equal(val, []byte("value_1")) {
		t.Errorf("Expected 'value_1', got '%s' (err: %v)", string(val), err)
	}

	err = database.Set("alpha", []byte("value_updated"))
	if err != nil {
		t.Fatalf("Failed to update alpha: %v", err)
	}

	valUpdated, err := database.Get("alpha")
	if err != nil || !bytes.Equal(valUpdated, []byte("value_updated")) {
		t.Errorf("Expected 'value_updated', got '%s'", string(valUpdated))
	}
}
