package command

import (
	"simple-database/pkg/database"
	"testing"
)

func TestCommandParser(t *testing.T) {

	tests := []struct {
		input    string
		expected CommandType
		args     []string
	}{
		{"SET key value", CmdSet, []string{"key", "value"}},
		{"GET key", CmdGet, []string{"key"}},
		{"UNSET key", CmdUnset, []string{"key"}},
		{"NUMEQUALTO value", CmdNumEqualTo, []string{"value"}},
		{"BEGIN", CmdBegin, []string{}},
		{"ROLLBACK", CmdRollback, []string{}},
		{"COMMIT", CmdCommit, []string{}},
		{"END", CmdEnd, []string{}},
		{"", CmdInvalid, nil},
		{"INVALID", CmdInvalid, nil},
		{"SET key", CmdInvalid, nil},     // Missing argument
		{"GET", CmdInvalid, nil},         // Missing argument
		{"BEGIN extra", CmdInvalid, nil}, // Extra argument
	}

	for _, test := range tests {
		cmd := parseCommand(test.input)
		if cmd.Type != test.expected {
			t.Errorf("Input '%s': expected %v, got %v", test.input, test.expected, cmd.Type)
		}

		if len(cmd.Args) != len(test.args) {
			t.Errorf("Input '%s': expected %d args, got %d", test.input, len(test.args), len(cmd.Args))
			continue
		}

		for i, arg := range test.args {
			if cmd.Args[i] != arg {
				t.Errorf("Input '%s': expected arg[%d] = '%s', got '%s'", test.input, i, arg, cmd.Args[i])
			}
		}
	}
}

func TestCommandExecutor(t *testing.T) {
	db := database.New()
	executor := NewExecutor(db)

	// Test SET command (no output)
	output, shouldExit := executor.Execute("SET test value")
	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}
	if shouldExit {
		t.Error("SET should not exit")
	}

	// Test GET command
	output, shouldExit = executor.Execute("GET test")
	if output != "value" {
		t.Errorf("Expected 'value', got '%s'", output)
	}
	if shouldExit {
		t.Error("GET should not exit")
	}

	// Test GET non-existent key
	output, shouldExit = executor.Execute("GET missing")

	if shouldExit {
		t.Error("GET should not exit")
	}

	if output != "NULL" {
		t.Errorf("Expected 'NULL', got '%s'", output)
	}

	// Test NUMEQUALTO
	output, shouldExit = executor.Execute("NUMEQUALTO value")
	if shouldExit {
		t.Error("NUMEQUALTO should not exit")
	}

	if output != "1" {
		t.Errorf("Expected '1', got '%s'", output)
	}

	// Test transaction commands
	output, shouldExit = executor.Execute("BEGIN")
	if shouldExit {
		t.Error("BEGIN should not exit")
	}
	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}

	output, shouldExit = executor.Execute("ROLLBACK")
	if shouldExit {
		t.Error("ROLLBACK should not exit")
	}
	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}

	// Test rollback with no transaction
	output, shouldExit = executor.Execute("ROLLBACK")
	if shouldExit {
		t.Error("ROLLBACK should not exit")
	}
	if output != "NO TRANSACTION" {
		t.Errorf("Expected 'NO TRANSACTION', got '%s'", output)
	}

	// Test commit with no transaction
	output, shouldExit = executor.Execute("COMMIT")
	if shouldExit {
		t.Error("COMMIT should not exit")
	}
	if output != "NO TRANSACTION" {
		t.Errorf("Expected 'NO TRANSACTION', got '%s'", output)
	}

	// Test END command
	output, shouldExit = executor.Execute("END")
	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}
	if !shouldExit {
		t.Error("END should exit")
	}

	// Test invalid command
	output, shouldExit = executor.Execute("INVALID")
	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}
	if shouldExit {
		t.Error("Invalid command should not exit")
	}
}

func TestCaseSensitivity(t *testing.T) {
	// Commands should be case-insensitive
	tests := []string{"set key value", "SET key value", "Set Key Value"}

	for _, test := range tests {
		cmd := parseCommand(test)
		if cmd.Type != CmdSet {
			t.Errorf("Input '%s': expected CmdSet, got %v", test, cmd.Type)
		}
	}

	// But keys and values should be preserved as-is
	cmd := parseCommand("SET MyKey MyValue")
	if cmd.Args[0] != "MyKey" {
		t.Errorf("Expected key 'MyKey', got '%s'", cmd.Args[0])
	}
	if cmd.Args[1] != "MyValue" {
		t.Errorf("Expected value 'MyValue', got '%s'", cmd.Args[1])
	}
}
