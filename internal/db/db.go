package db

import (
	"dbdb/internal/storage"
	"dbdb/internal/tree"
	"sync"
)

type DB struct {
	store *storage.Physical
	tree  *tree.BTree
	mu    sync.RWMutex
}

func Open(path string) (*DB, error) {
	p, err := storage.OpenPhysical(path)
	if err != nil {
		return nil, err
	}

	t := tree.NewBTree(p, tree.NodeRef{})

	return &DB{
		store: p,
		tree:  t,
	}, nil
}

func (db *DB) Get(key string) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.tree.Get(key)
}

func (db *DB) Set(key string, value []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	newRoot, err := db.tree.Set(key, value)
	if err != nil {
		return err
	}

	db.tree.Root = newRoot
	return nil
}
