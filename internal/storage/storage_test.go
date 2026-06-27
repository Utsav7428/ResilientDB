package storage_test

import (
	"bytes"
	"dbdb/internal/storage"
	"os"
	"sync"
	"testing"
)

func setupTempFile(t *testing.T) string {
	file, err := os.CreateTemp("", "dbdb_storage_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	file.Close()
	t.Cleanup(func() { os.Remove(file.Name()) })
	return file.Name()
}

func TestPhysicalAppendAndRead(t *testing.T) {
	path := setupTempFile(t)
	store, err := storage.OpenPhysical(path)
	if err != nil {
		t.Fatalf("Failed to open physical store: %v", err)
	}

	data1 := []byte("hello")
	data2 := []byte("world")

	offset1, err := store.Append(data1)
	if err != nil || offset1 != 0 {
		t.Fatalf("Expected offset 0, got %d (err: %v)", offset1, err)
	}

	offset2, err := store.Append(data2)
	if err != nil || offset2 != int64(len(data1)) {
		t.Fatalf("Expected offset %d, got %d", len(data1), offset2)
	}

	read2, err := store.Read(offset2, len(data2))
	if err != nil {
		t.Fatalf("Failed to read back data: %v", err)
	}
	if !bytes.Equal(read2, data2) {
		t.Errorf("Expected %s, got %s", data2, read2)
	}
}

func TestConcurrentAppends(t *testing.T) {
	path := setupTempFile(t)
	store, err := storage.OpenPhysical(path)
	if err != nil {
		t.Fatalf("Failed to open physical store: %v", err)
	}

	var wg sync.WaitGroup
	routines := 50
	payload := []byte("concurrent_write")

	errChan := make(chan error, routines)

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Append(payload)
			if err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("Concurrent write failed: %v", err)
	}

	// Verify final file size (50 writes * 16 bytes)
	stat, _ := os.Stat(path)
	expectedSize := int64(routines * len(payload))
	if stat.Size() != expectedSize {
		t.Errorf("Expected file size %d, got %d. Mutex failed to lock correctly.", expectedSize, stat.Size())
	}
}
