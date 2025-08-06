package main

import (
	"bufio"
	"fmt"
	"os"
	"simple-database/pkg/command"
	"simple-database/pkg/database"
)

func main() {
	db := database.New()
	executor := command.NewExecutor(db)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		shouldExit := executor.ExecuteAndPrint(input)

		if shouldExit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
}
