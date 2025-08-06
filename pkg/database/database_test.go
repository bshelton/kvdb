package database

import (
	"testing"
)

func TestBasicOperations(t *testing.T) {
	db := New()
	db.Set("key1", "value1")
	if got := db.Get("key1"); got != "value1" {
		t.Errorf("Expected 'value1', got '%s'", got)
	}

	if got := db.Get("nonexistent"); got != "NULL" {
		t.Errorf("Expected 'NULL', got '%s'", got)
	}

	db.Unset("key1")
	if got := db.Get("key1"); got != "NULL" {
		t.Errorf("Expected 'NULL' after unset, got '%s'", got)
	}
}

func TestScenario1(t *testing.T) {
	db := New()
	db.Set("ex", "10")

	if got := db.Get("ex"); got != "10" {
		t.Errorf("Expected '10', got '%s'", got)
	}

	db.Unset("ex")

	if got := db.Get("EX"); got != "NULL" {
		t.Errorf("Expected 'NULL', got '%s'", got)
	}
}

func TestScenario2(t *testing.T) {
	db := New()
	db.Set("a", "10")
	db.Set("b", "10")

	if got := db.NumEqualTo("10"); got != 2 {
		t.Errorf("Expected 2, got %d", got)
	}

	if got := db.NumEqualTo("20"); got != 0 {
		t.Errorf("Expected 0, got %d", got)
	}

	db.Set("b", "30")

	if got := db.NumEqualTo("10"); got != 1 {
		t.Errorf("Expected 1, got %d", got)
	}
}

func TestTransactionScenario1(t *testing.T) {
	db := New()
	db.Set("a", "10")
	if got := db.Get("a"); got != "10" {
		t.Errorf("Expected '10', got '%s'", got)
	}

	db.Begin()
	db.Set("a", "20")

	if got := db.Get("a"); got != "20" {
		t.Errorf("Expected '20', got '%s'", got)
	}

	if err := db.Rollback(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := db.Get("a"); got != "10" {
		t.Errorf("Expected '10', got '%s'", got)
	}

	if err := db.Rollback(); err == nil {
		t.Error("Expected error for rollback with no transaction")
	}

	if got := db.Get("a"); got != "10" {
		t.Errorf("Expected '10', got '%s'", got)
	}
}

func TestTransactionScenario2(t *testing.T) {
	db := New()
	db.Begin()
	db.Set("a", "30")
	db.Begin()
	db.Set("a", "40")
	if err := db.Commit(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := db.Get("a"); got != "40" {
		t.Errorf("Expected '40', got '%s'", got)
	}

	if err := db.Rollback(); err == nil {
		t.Error("Expected error for rollback with no transaction")
	}
}

func TestTransactionScenario3(t *testing.T) {
	db := New()

	db.Set("a", "50")

	db.Begin()

	if got := db.Get("a"); got != "50" {
		t.Errorf("Expected '50', got '%s'", got)
	}

	db.Set("a", "60")
	db.Begin()
	db.Unset("a")

	if got := db.Get("a"); got != "NULL" {
		t.Errorf("Expected 'NULL', got '%s'", got)
	}

	if err := db.Rollback(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := db.Get("a"); got != "60" {
		t.Errorf("Expected '60', got '%s'", got)
	}

	if err := db.Commit(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := db.Get("a"); got != "60" {
		t.Errorf("Expected '60', got '%s'", got)
	}
}

func TestTransactionScenario4(t *testing.T) {
	db := New()
	db.Set("a", "10")
	db.Begin()

	if got := db.NumEqualTo("10"); got != 1 {
		t.Errorf("Expected 1, got %d", got)
	}

	db.Begin()
	db.Unset("a")

	if got := db.NumEqualTo("10"); got != 0 {
		t.Errorf("Expected 0, got %d", got)
	}

	if err := db.Rollback(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := db.NumEqualTo("10"); got != 1 {
		t.Errorf("Expected 1, got %d", got)
	}

	if err := db.Commit(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestNestedTransactions(t *testing.T) {
	db := New()

	// Test deeply nested transactions
	db.Set("key", "original")

	db.Begin()
	db.Set("key", "level1")

	db.Begin()
	db.Set("key", "level2")

	db.Begin()
	db.Set("key", "level3")

	// Should see level3 value
	if got := db.Get("key"); got != "level3" {
		t.Errorf("Expected 'level3', got '%s'", got)
	}

	// Rollback one level
	if err := db.Rollback(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := db.Get("key"); got != "level2" {
		t.Errorf("Expected 'level2', got '%s'", got)
	}

	// Commit should apply all remaining changes
	if err := db.Commit(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := db.Get("key"); got != "level2" {
		t.Errorf("Expected 'level2', got '%s'", got)
	}
}

func TestValueCounting(t *testing.T) {
	db := New()

	// Test value counting with various operations
	db.Set("a", "100")
	db.Set("b", "100")
	db.Set("c", "200")

	if got := db.NumEqualTo("100"); got != 2 {
		t.Errorf("Expected 2, got %d", got)
	}

	if got := db.NumEqualTo("200"); got != 1 {
		t.Errorf("Expected 1, got %d", got)
	}

	db.Set("a", "200")

	if got := db.NumEqualTo("100"); got != 1 {
		t.Errorf("Expected 1, got %d", got)
	}

	if got := db.NumEqualTo("200"); got != 2 {
		t.Errorf("Expected 2, got %d", got)
	}

	db.Unset("b")

	if got := db.NumEqualTo("100"); got != 0 {
		t.Errorf("Expected 0, got %d", got)
	}
}

func TestTransactionValueCounting(t *testing.T) {
	db := New()

	db.Set("a", "100")
	db.Set("b", "100")

	if got := db.NumEqualTo("100"); got != 2 {
		t.Errorf("Expected 2, got %d", got)
	}

	db.Begin()
	db.Set("c", "100")

	if got := db.NumEqualTo("100"); got != 3 {
		t.Errorf("Expected 3, got %d", got)
	}

	db.Unset("a")

	if got := db.NumEqualTo("100"); got != 2 {
		t.Errorf("Expected 2, got %d", got)
	}

	if err := db.Rollback(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := db.NumEqualTo("100"); got != 2 {
		t.Errorf("Expected 2, got %d", got)
	}
}

func TestErrorHandling(t *testing.T) {
	db := New()

	// Test rollback with no transaction
	if err := db.Rollback(); err == nil {
		t.Error("Expected error for rollback with no transaction")
	}

	if err := db.Rollback(); err.Error() != "NO TRANSACTION" {
		t.Errorf("Expected 'NO TRANSACTION', got '%s'", err.Error())
	}

	// Test commit with no transaction
	if err := db.Commit(); err == nil {
		t.Error("Expected error for commit with no transaction")
	}

	if err := db.Commit(); err.Error() != "NO TRANSACTION" {
		t.Errorf("Expected 'NO TRANSACTION', got '%s'", err.Error())
	}
}
