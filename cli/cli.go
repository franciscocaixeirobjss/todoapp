package cli

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"todoapp/task"
)

func RunCLI() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Interactive CLI started. Type 'help' for available commands.")
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		switch {
		case command == "get":
			printTasks()
		case strings.HasPrefix(command, "create"):
			handleCreate(command)
		case command == "exit":
			fmt.Println("Exiting CLI...")
			return
		case command == "help":
			fmt.Println("Available commands:")
			fmt.Println("  get                                   - Retrieve and display all tasks")
			fmt.Println("  create -title <title> -desc <description> -status <status> - Create a new task with the given details.")
			fmt.Println("      Example: create -title \"Golang\" -desc \"Task1\" -status \"NotStarted\"")
			fmt.Println("  exit                                  - Exit the CLI")
			fmt.Println("  help                                  - Show this help message")
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func printTasks() {
	tasks := task.GetTasks()
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	tasksJSON, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal tasks", "error", err)
		return
	}
	fmt.Println(string(tasksJSON))
}

func handleCreate(command string) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	title := createCmd.String("title", "", "Title of the task")
	description := createCmd.String("description", "", "Description of the task")
	status := createCmd.String("status", "", "Status of the task (NotStarted, Started, Completed)")

	args := strings.Fields(command)
	if len(args) < 2 {
		fmt.Println("Usage: create -title <title> -description <description> -status <status>")
		return
	}
	err := createCmd.Parse(args[1:])
	if err != nil {
		fmt.Println("Failed to parse arguments:", err)
		return
	}

	newTask := task.Task{
		Title:        *title,
		Description:  *description,
		StatusString: *status,
	}
	err = task.CreateTask(newTask)
	if err != nil {
		fmt.Println("Failed to create task:", err)
		return
	}
	fmt.Println("Task created successfully.")
}
