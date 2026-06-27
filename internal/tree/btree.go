package tree

import (
	"bytes"
	"dbdb/internal/storage"
	"encoding/gob"
	"errors"
)

type NodeRef struct {
	Offset int64
	Length int
}

type Node struct {
	Key   string
	Value []byte
	Left  NodeRef
	Right NodeRef
}

type BTree struct {
	store *storage.Physical
	Root  NodeRef
}

func NewBTree(s *storage.Physical, root NodeRef) *BTree {
	return &BTree{
		store: s,
		Root:  root,
	}
}

// readNode deserializes a node block directly from its disk offset.
func (t *BTree) readNode(ref NodeRef) (*Node, error) {
	if ref.Length == 0 {
		return nil, nil
	}

	data, err := t.store.Read(ref.Offset, ref.Length)
	if err != nil {
		return nil, err
	}

	var node Node
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&node); err != nil {
		return nil, err
	}
	return &node, nil
}

// writeNode serializes a node and appends it cleanly to the end of the file.
func (t *BTree) writeNode(node *Node) (NodeRef, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(node); err != nil {
		return NodeRef{}, err
	}

	data := buf.Bytes()
	offset, err := t.store.Append(data)
	if err != nil {
		return NodeRef{}, err
	}

	return NodeRef{Offset: offset, Length: len(data)}, nil
}

// Get searches for a key recursively through the immutable disk nodes.
func (t *BTree) Get(key string) ([]byte, error) {
	return t.get(t.Root, key)
}

func (t *BTree) get(ref NodeRef, key string) ([]byte, error) {
	if ref.Length == 0 {
		return nil, errors.New("key not found")
	}

	node, err := t.readNode(ref)
	if err != nil {
		return nil, err
	}

	if key == node.Key {
		return node.Value, nil
	} else if key < node.Key {
		return t.get(node.Left, key)
	} else {
		return t.get(node.Right, key)
	}
}

// Set performs an immutable insert, returning a brand new updated Root reference.
func (t *BTree) Set(key string, value []byte) (NodeRef, error) {
	return t.set(t.Root, key, value)
}

func (t *BTree) set(ref NodeRef, key string, value []byte) (NodeRef, error) {
	if ref.Length == 0 {
		return t.writeNode(&Node{Key: key, Value: value})
	}

	node, err := t.readNode(ref)
	if err != nil {
		return NodeRef{}, err
	}

	var nextNode Node
	nextNode.Key = node.Key
	nextNode.Value = node.Value
	nextNode.Left = node.Left
	nextNode.Right = node.Right

	if key == node.Key {
		nextNode.Value = value
	} else if key < node.Key {
		newLeft, err := t.set(node.Left, key, value)
		if err != nil {
			return NodeRef{}, err
		}
		nextNode.Left = newLeft
	} else {
		newRight, err := t.set(node.Right, key, value)
		if err != nil {
			return NodeRef{}, err
		}
		nextNode.Right = newRight
	}

	return t.writeNode(&nextNode)
}
