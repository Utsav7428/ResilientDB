# Res Key-Value Store

DBDB is an on-disk, key-value storage engine modeled as a binary search tree. It is built to safely store, update, and retrieve data directly from your computer's hard drive.

## Core Capability

The main advantage of this engine is crash resilience through immutable architecture.

Traditional databases modify files by overwriting data in-place. If the power cuts out or the system crashes mid-write, the data becomes corrupted. DBDB solves this by treating the storage file as append-only. 

When you insert or update a key, the engine never modifies old data. Instead, it writes a brand-new copy of the updated branch to the very end of the file and moves its master pointer to this new location only after the write is fully completed. If a crash happens midway, the database simply falls back to the last known complete state, making data corruption virtually impossible.



## Project Structure

* cmd/dbdb/main.go: Entry point to interact with the key-value engine.
* internal/storage/physical.go: Low-level file coordinator that handles safe, concurrent disk writing and reading.
* internal/tree/btree.go: Manages the logical sorting and searching paths of the tree layout.
* internal/db/db.go: High-level database manager that enforces transactional safety locks for reads and writes.

## Automation Scripts

* build.bat: Compiles the source files into a single executable binary in the bin directory.
* test.bat: Runs all unit and integration tests, including simulated crash endurance and concurrent writing stresses.
