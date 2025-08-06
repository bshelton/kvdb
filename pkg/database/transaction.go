package database

import "errors"

var (
	ErrNoTransaction = errors.New("NO TRANSACTION")
)

// Operation represents the type of transaction operation
type Operation int

const (
	OpSet Operation = iota
	OpUnset
)

// TransactionChange represents a change made within a transaction
type TransactionChange struct {
	Key       string
	OldValue  string
	NewValue  string
	Operation Operation
}

// TransactionLayer represents a single transaction layer
type TransactionLayer struct {
	changes     map[string]TransactionChange
	valueCounts map[string]int
}

// newTransactionLayer creates a new transaction layer
func newTransactionLayer() *TransactionLayer {
	return &TransactionLayer{
		changes:     make(map[string]TransactionChange),
		valueCounts: make(map[string]int),
	}
}

// TransactionManager manages nested transactions
type TransactionManager struct {
	layers []*TransactionLayer
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		layers: make([]*TransactionLayer, 0),
	}
}

// InTransaction returns true if there are active transactions
func (tm *TransactionManager) InTransaction() bool {
	return len(tm.layers) > 0
}

// Begin starts a new transaction layer
func (tm *TransactionManager) Begin() {
	tm.layers = append(tm.layers, newTransactionLayer())
}

// Rollback removes the most recent transaction layer
func (tm *TransactionManager) Rollback() error {
	if !tm.InTransaction() {
		return ErrNoTransaction
	}
	tm.layers = tm.layers[:len(tm.layers)-1]
	return nil
}

// Set records a SET operation in the current transaction
func (tm *TransactionManager) Set(key, value, oldValue string) {
	if !tm.InTransaction() {
		return
	}

	layer := tm.getCurrentLayer()

	if oldValue != "NULL" {
		layer.valueCounts[oldValue]--
	}

	layer.valueCounts[value]++
	layer.changes[key] = TransactionChange{
		Key:       key,
		OldValue:  oldValue,
		NewValue:  value,
		Operation: OpSet,
	}
}

// Unset records an UNSET operation in the current transaction
func (tm *TransactionManager) Unset(key, currentValue string) {
	if !tm.InTransaction() {
		return
	}

	layer := tm.getCurrentLayer()
	layer.valueCounts[currentValue]--
	layer.changes[key] = TransactionChange{
		Key:       key,
		OldValue:  currentValue,
		NewValue:  "NULL",
		Operation: OpUnset,
	}
}

// Get retrieves a value from the transaction layers
func (tm *TransactionManager) Get(key string) (string, bool) {
	// Search from most recent transaction to oldest
	for i := len(tm.layers) - 1; i >= 0; i-- {
		if change, exists := tm.layers[i].changes[key]; exists {
			if change.Operation == OpUnset {
				return "NULL", true
			}
			return change.NewValue, true
		}
	}
	return "", false
}

// GetValueCount returns the net change in value count across all transaction layers
func (tm *TransactionManager) GetValueCount(value string) int {
	total := 0
	for _, layer := range tm.layers {
		total += layer.valueCounts[value]
	}
	return total
}

// GetAllChanges returns all changes from all transaction layers
func (tm *TransactionManager) GetAllChanges() []TransactionChange {
	var allChanges []TransactionChange

	// Collect changes from all layers, with later layers overriding earlier ones
	changeMap := make(map[string]TransactionChange)

	for _, layer := range tm.layers {
		for key, change := range layer.changes {
			changeMap[key] = change
		}
	}

	for _, change := range changeMap {
		allChanges = append(allChanges, change)
	}

	return allChanges
}

// Clear removes all transaction layers
func (tm *TransactionManager) Clear() {
	tm.layers = make([]*TransactionLayer, 0)
}

// getCurrentLayer returns the current (most recent) transaction layer
func (tm *TransactionManager) getCurrentLayer() *TransactionLayer {
	if !tm.InTransaction() {
		return nil
	}
	return tm.layers[len(tm.layers)-1]
}
