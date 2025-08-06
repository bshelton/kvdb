package command

import (
	"fmt"
	"strconv"
	"strings"
)

// CommandType represents different database commands
type CommandType int

const (
	CmdSet CommandType = iota
	CmdGet
	CmdUnset
	CmdNumEqualTo
	CmdBegin
	CmdRollback
	CmdCommit
	CmdEnd
	CmdInvalid
)

// Command represents a parsed database command
type Command struct {
	Type CommandType
	Args []string
}

// parseCommand converts a string input into a Command
func parseCommand(input string) Command {
	input = strings.TrimSpace(input)
	if input == "" {
		return Command{Type: CmdInvalid}
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return Command{Type: CmdInvalid}
	}

	cmdName := strings.ToUpper(parts[0])
	args := parts[1:]

	switch cmdName {
	case "SET":
		if len(args) == 2 {
			return Command{Type: CmdSet, Args: args}
		}
	case "GET":
		if len(args) == 1 {
			return Command{Type: CmdGet, Args: args}
		}
	case "UNSET":
		if len(args) == 1 {
			return Command{Type: CmdUnset, Args: args}
		}
	case "NUMEQUALTO":
		if len(args) == 1 {
			return Command{Type: CmdNumEqualTo, Args: args}
		}
	case "BEGIN":
		if len(args) == 0 {
			return Command{Type: CmdBegin}
		}
	case "ROLLBACK":
		if len(args) == 0 {
			return Command{Type: CmdRollback}
		}
	case "COMMIT":
		if len(args) == 0 {
			return Command{Type: CmdCommit}
		}
	case "END":
		if len(args) == 0 {
			return Command{Type: CmdEnd}
		}
	}
	return Command{Type: CmdInvalid}
}

// Database defines the interface that the command executor expects
type Database interface {
	Set(key, value string)
	Get(key string) string
	Unset(key string)
	NumEqualTo(value string) int
	Begin()
	Rollback() error
	Commit() error
}

// Executor handles the execution of database commands
type Executor struct {
	database Database
}

// NewExecutor creates a new command executor
func NewExecutor(db Database) *Executor {
	return &Executor{
		database: db,
	}
}

// Execute processes a command string and returns output if any
func (ce *Executor) Execute(input string) (output string, shouldExit bool) {
	cmd := parseCommand(input)

	switch cmd.Type {
	case CmdSet:
		ce.database.Set(cmd.Args[0], cmd.Args[1])
		return "", false

	case CmdGet:
		result := ce.database.Get(cmd.Args[0])
		return result, false

	case CmdUnset:
		ce.database.Unset(cmd.Args[0])
		return "", false

	case CmdNumEqualTo:
		count := ce.database.NumEqualTo(cmd.Args[0])
		return strconv.Itoa(count), false

	case CmdBegin:
		ce.database.Begin()
		return "", false

	case CmdRollback:
		if err := ce.database.Rollback(); err != nil {
			return err.Error(), false
		}
		return "", false

	case CmdCommit:
		if err := ce.database.Commit(); err != nil {
			return err.Error(), false
		}
		return "", false

	case CmdEnd:
		return "", true

	case CmdInvalid:
		return "", false
	}

	return "", false
}

// ExecuteAndPrint processes a command and prints output if needed
func (ce *Executor) ExecuteAndPrint(input string) bool {
	output, shouldExit := ce.Execute(input)
	if output != "" {
		fmt.Println(output)
	}
	return shouldExit
}
