package database

import "simple-database/pkg/storage"

// Database represents an in-memory key-value store with transaction support
type Database struct {
	storage      *storage.Storage
	transactions *TransactionManager
}

// New creates a new database instance
func New() *Database {
	return &Database{
		storage:      storage.New(),
		transactions: NewTransactionManager(),
	}
}

// Set stores a key-value pair
func (db *Database) Set(key, value string) {
	if db.transactions.InTransaction() {
		db.transactions.Set(key, value, db.storage.Get(key))
	} else {
		oldValue := db.storage.Get(key)
		db.storage.Set(key, value)
		db.storage.UpdateValueCount(oldValue, value)
	}
}

// Get retrieves a value by key, returns "NULL" if not found
func (db *Database) Get(key string) string {
	if db.transactions.InTransaction() {
		if value, found := db.transactions.Get(key); found {
			return value
		}
	}
	return db.storage.Get(key)
}

// Unset removes a key-value pair
func (db *Database) Unset(key string) {
	currentValue := db.Get(key)
	if currentValue == "NULL" {
		return
	}

	if db.transactions.InTransaction() {
		db.transactions.Unset(key, currentValue)
	} else {
		db.storage.Unset(key)
		db.storage.DecrementValueCount(currentValue)
	}
}

// NumEqualTo returns the count of keys with the given value
func (db *Database) NumEqualTo(value string) int {
	baseCount := db.storage.GetValueCount(value)
	transactionCount := db.transactions.GetValueCount(value)
	return baseCount + transactionCount
}

// Begin starts a new transaction
func (db *Database) Begin() {
	db.transactions.Begin()
}

// Rollback undoes the most recent transaction
func (db *Database) Rollback() error {
	return db.transactions.Rollback()
}

// Commit applies all pending transactions to the main storage
func (db *Database) Commit() error {
	if !db.transactions.InTransaction() {
		return ErrNoTransaction
	}

	changes := db.transactions.GetAllChanges()
	for _, change := range changes {
		db.applyChange(change)
	}

	db.transactions.Clear()
	return nil
}

// applyChange applies a single transaction change to the main storage
func (db *Database) applyChange(change TransactionChange) {
	switch change.Operation {
	case OpSet:
		oldValue := db.storage.Get(change.Key)
		db.storage.Set(change.Key, change.NewValue)
		db.storage.UpdateValueCount(oldValue, change.NewValue)
	case OpUnset:
		db.storage.Unset(change.Key)
		db.storage.DecrementValueCount(change.OldValue)
	}
}
