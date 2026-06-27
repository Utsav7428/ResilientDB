package main

import (
	"fmt"
	"dbdb/internal/storage"
)

func main() {
	fmt.Println("Initializing DBDB Key-Value Engine...")
	_ = &storage.Physical{}
}
