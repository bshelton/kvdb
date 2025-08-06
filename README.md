# Simple Database

This is my implementation of a simple in-memory key-value database with transaction support.

## What It Does

- Store and retrieve key-value pairs
- Count how many keys have a particular value
- Full transaction support with nesting
- All data is kept in memory (nothing gets saved to disk)

## Requirements

```bash
# go version 1.23
go mod tidy
go run main.go
```

Once it's running, you can type commands directly or pipe them in from a file.

## Running the Examples

I've included several example files in the `examples/` folder that demonstrate different features.

```bash
# Basic operations (SET, GET, UNSET, case sensitivity)
go run main.go < examples/test_scenario1.txt

# Value counting with NUMEQUALTO
go run main.go < examples/test_scenario2.txt

# Transaction examples
go run main.go < examples/test_transaction1.txt   # Basic transaction with rollback
go run main.go < examples/test_transaction2.txt   # Nested transactions with commit
go run main.go < examples/test_transaction3.txt   # Complex nested rollback scenario
go run main.go < examples/test_transaction4.txt   # Transaction value counting
```

Each example file contains a series of commands that will be executed in sequence, and you'll see the output for commands that produce results (like GET and NUMEQUALTO).

## Commands

### Basic Data Operations

- `SET key value` - Store a value with a key
- `GET key` - Get the value for a key (returns "NULL" if not found)
- `UNSET key` - Remove a key and its value
- `NUMEQUALTO value` - Count how many keys currently have this value

### Transaction Commands

- `BEGIN` - Start a new transaction (you can nest these)
- `ROLLBACK` - Undo everything in the most recent transaction
- `COMMIT` - Apply all pending transaction changes permanently

### Control

- `END` - Exit the program

## Examples

Here's what typical usage looks like:

```
SET name Alice
GET name
Alice

SET age 30
SET name Bob
GET name
Bob

NUMEQUALTO Bob
1

UNSET name
GET name
NULL
```

And here's how transactions work:

```
SET balance 100
BEGIN
SET balance 50
GET balance
50
ROLLBACK
GET balance
100
```

You can even nest transactions:

```
SET x 1
BEGIN
SET x 2
BEGIN
SET x 3
ROLLBACK
GET x
2
COMMIT
GET x
2
```

## Design Decisions & Assumptions

### Transaction Behavior

- **Isolation**: Changes inside transactions are isolated until you commit them
- **Nesting**: You can have transactions inside transactions. ROLLBACK undoes just the innermost one, but COMMIT applies everything
- **Error Handling**: If you try to ROLLBACK or COMMIT without an active transaction, you get "NO TRANSACTION"

### Value Counting

- Track value counts separately from the main storage to make NUMEQUALTO fast (O(1) instead of scanning everything)
- Transaction changes update these counts incrementally, so the counts stay accurate even with uncommitted changes

### Key/Value Rules

- **Case Sensitivity**: Keys are case-sensitive ("key" and "KEY" are different)
- **String Only**: Everything is stored as strings
- **Command Case**: Commands themselves are case-insensitive (SET, set, Set all work)

### Input Handling

- The program reads from stdin line by line
- It handles EOF gracefully (like when you pipe in a file or press Ctrl+D)
- Invalid commands are ignored silently (this matched the behavior described in the original challenge)
- Empty lines get skipped

### Memory Management

- Everything lives in memory using Go's built-in maps
- When you UNSET something or a value count drops to zero, memory is cleaned up
- No data persists between program runs

## How It's Organized

The code is now organized into clean packages:

- `main.go` - Entry point that coordinates everything
- `pkg/database/` - Core database logic and transaction management
  - `database.go` - Main database interface
  - `transaction.go` - Transaction management system
  - `database_test.go` - Database and transaction tests
- `pkg/storage/` - Key-value storage and counting
  - `storage.go` - Core storage operations
- `pkg/command/` - Command parsing and execution
  - `command.go` - Command parser and executor
  - `command_test.go` - Command parsing and execution tests

The transaction system was the most interesting challenge. I used a stack of "layers" where each BEGIN adds a new layer, and changes get recorded there. ROLLBACK just throws away the top layer, while COMMIT merges all layers down into the main storage.

## Testing

The tests are now organized alongside their respective code in each package:

- **Database Tests** (`pkg/database/database_test.go`)

  - All original PDF test scenarios (scenarios 1-4, transactions 1-4)
  - Edge cases like deeply nested transactions
  - Error conditions and boundary cases
  - Transaction isolation and rollback behavior

- **Command Tests** (`pkg/command/command_test.go`)
  - Command parsing with various inputs
  - Command execution and output validation
  - Case sensitivity handling
  - Invalid input handling

Run all tests with:

```bash
go test ./...     # Run tests in all packages
```

Or run tests for specific packages:

```bash
go test ./pkg/database -v    # Database and transaction tests
go test ./pkg/command -v     # Command parsing tests
go test ./pkg/storage -v     # Storage tests (if any)
```

## Limitations

A few things this doesn't do (by design):

- No disk persistence - everything disappears when you exit
- No networking or concurrent access
- Only handles string values
- Memory usage grows with your data (no automatic cleanup)

## Why I Built It This Way

I focused on clean, readable code over micro-optimizations. I used Go's built-in data structures rather than implementing custom ones because they're well-tested and performant enough for this use case.

The modular structure makes it easy to extend - for example, you could add new commands by just updating the command parser and executor, or add persistence by implementing a new storage backend.
