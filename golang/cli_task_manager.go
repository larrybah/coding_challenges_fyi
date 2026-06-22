package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Task represents a single task in the task manager
// Struct: a collection of fields that group related data
type Task struct {
	ID        int       `json:"id"`        // Unique identifier for the task
	Title     string    `json:"title"`     // The task description/title
	Completed bool      `json:"completed"` // Whether the task is done or not
	CreatedAt time.Time `json:"created_at"` // When the task was created
	UpdatedAt time.Time `json:"updated_at"` // When the task was last modified
}

// TaskManager holds the collection of tasks and provides methods to manage them
// This is the main structure for our application
type TaskManager struct {
	Tasks    []Task // Slice of Task structs
	FilePath string // Path where tasks are saved to disk
}

// NewTaskManager creates a new TaskManager instance and loads existing tasks from file
// Constructor pattern: initializes and returns a new instance
func NewTaskManager(filePath string) *TaskManager {
	tm := &TaskManager{
		Tasks:    []Task{},
		FilePath: filePath,
	}
	// Try to load existing tasks from file (if file exists)
	tm.LoadTasks()
	return tm
}

// LoadTasks reads tasks from a JSON file and populates the Tasks slice
// File Handling: reading and parsing JSON data from disk
func (tm *TaskManager) LoadTasks() error {
	// Check if file exists
	if _, err := os.Stat(tm.FilePath); os.IsNotExist(err) {
		// File doesn't exist yet, start with empty tasks
		fmt.Printf("No existing tasks file found. Creating new task manager.\n\n")
		return nil
	}

	// Read the entire file into memory
	data, err := ioutil.ReadFile(tm.FilePath)
	if err != nil {
		fmt.Printf("Error reading tasks file: %v\n", err)
		return err
	}

	// Check if file is empty
	if len(data) == 0 {
		return nil
	}

	// JSON Decoding: parse JSON string into Go struct
	// json.Unmarshal converts JSON bytes into our Task struct slice
	err = json.Unmarshal(data, &tm.Tasks)
	if err != nil {
		fmt.Printf("Error parsing tasks file: %v\n", err)
		return err
	}

	fmt.Printf("Loaded %d tasks from file.\n\n", len(tm.Tasks))
	return nil
}

// SaveTasks writes all tasks to the JSON file
// File Handling: writing and formatting JSON data to disk
func (tm *TaskManager) SaveTasks() error {
	// JSON Encoding: convert Go struct slice into JSON bytes
	// MarshalIndent adds nice formatting (4-space indent) for readability
	data, err := json.MarshalIndent(tm.Tasks, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling tasks: %v\n", err)
		return err
	}

	// Write the JSON data to file
	// 0644 permissions: owner can read/write, others can only read
	err = ioutil.WriteFile(tm.FilePath, data, 0644)
	if err != nil {
		fmt.Printf("Error saving tasks file: %v\n", err)
		return err
	}

	return nil
}

// AddTask creates a new task and adds it to the manager
// Method: a function associated with a specific type (TaskManager)
func (tm *TaskManager) AddTask(title string) {
	// Find the next available ID by checking the highest current ID
	nextID := 1
	if len(tm.Tasks) > 0 {
		nextID = tm.Tasks[len(tm.Tasks)-1].ID + 1
	}

	// Create new Task struct instance
	task := Task{
		ID:        nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Append to the Tasks slice
	tm.Tasks = append(tm.Tasks, task)
	fmt.Printf("✓ Task added: %s (ID: %d)\n", title, nextID)

	// Persist to file
	tm.SaveTasks()
}

// CompleteTask marks a task as completed
// Takes a task ID and updates the Completed field
func (tm *TaskManager) CompleteTask(id int) {
	// Loop through all tasks to find the one with matching ID
	for i, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks[i].Completed = true
			tm.Tasks[i].UpdatedAt = time.Now()
			fmt.Printf("✓ Task completed: %s\n", task.Title)
			tm.SaveTasks()
			return
		}
	}
	fmt.Printf("✗ Task with ID %d not found\n", id)
}

// DeleteTask removes a task by ID
// Rebuilds the Tasks slice without the deleted task
func (tm *TaskManager) DeleteTask(id int) {
	// Create a new slice to hold tasks after deletion
	var newTasks []Task

	// Loop through tasks and keep only those that don't match the ID to delete
	for _, task := range tm.Tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		} else {
			fmt.Printf("✓ Task deleted: %s\n", task.Title)
		}
	}

	// If no task was deleted, notify the user
	if len(newTasks) == len(tm.Tasks) {
		fmt.Printf("✗ Task with ID %d not found\n", id)
		return
	}

	// Update the manager's tasks
	tm.Tasks = newTasks
	tm.SaveTasks()
}

// ListTasks displays all tasks in a formatted table
// Shows task ID, status, title, and creation date
func (tm *TaskManager) ListTasks() {
	if len(tm.Tasks) == 0 {
		fmt.Println("No tasks found. Add a task to get started!")
		return
	}

	fmt.Println("\n=== YOUR TASKS ===")
	fmt.Println(strings.Repeat("-", 60))
	// Print header row
	fmt.Printf("%-4s %-6s %-40s %-10s\n", "ID", "Status", "Title", "Created")
	fmt.Println(strings.Repeat("-", 60))

	// Print each task as a row
	for _, task := range tm.Tasks {
		// Display completed status as [✓] or [ ]
		status := "[ ]"
		if task.Completed {
			status = "[✓]"
		}

		// Format the creation date as MM/DD/YYYY
		createdDate := task.CreatedAt.Format("01/02/2006")

		// Truncate long titles to fit in the column
		title := task.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		fmt.Printf("%-4d %-6s %-40s %-10s\n", task.ID, status, title, createdDate)
	}

	fmt.Println(strings.Repeat("-", 60))
}

// SearchTasks finds and displays tasks matching a keyword
// Case-insensitive search through all task titles
func (tm *TaskManager) SearchTasks(keyword string) {
	var results []Task

	// Loop through all tasks and check if keyword is in the title
	for _, task := range tm.Tasks {
		// strings.Contains checks if substring exists (case-sensitive)
		// We convert both to lowercase for case-insensitive comparison
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(keyword)) {
			results = append(results, task)
		}
	}

	if len(results) == 0 {
		fmt.Printf("No tasks found matching '%s'\n", keyword)
		return
	}

	fmt.Printf("\n=== SEARCH RESULTS FOR '%s' ===\n", keyword)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-4s %-6s %-40s\n", "ID", "Status", "Title")
	fmt.Println(strings.Repeat("-", 60))

	for _, task := range results {
		status := "[ ]"
		if task.Completed {
			status = "[✓]"
		}

		title := task.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		fmt.Printf("%-4d %-6s %-40s\n", task.ID, status, title)
	}

	fmt.Println(strings.Repeat("-", 60))
}

// InteractiveMode starts an interactive CLI where users can input commands
// Reads from standard input and processes user commands
func (tm *TaskManager) InteractiveMode() {
	// Create a scanner to read from standard input (keyboard)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n=== TASK MANAGER ===")
	fmt.Println("Commands: add, list, complete, delete, search, quit")
	fmt.Println(strings.Repeat("-", 40))

	for {
		// Print prompt and get user input
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		// Get the input and split into command and arguments
		input := scanner.Text()
		parts := strings.Fields(input) // Split by whitespace

		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		// Handle different commands
		switch command {
		case "add":
			// "add" command requires a title argument
			if len(parts) < 2 {
				fmt.Println("Usage: add <task title>")
				continue
			}
			// Join all words after "add" to form the full title
			title := strings.Join(parts[1:], " ")
			tm.AddTask(title)

		case "list":
			tm.ListTasks()

		case "complete":
			// "complete" command requires a task ID
			if len(parts) < 2 {
				fmt.Println("Usage: complete <task_id>")
				continue
			}
			// Parse the ID string into an integer
			var id int
			_, err := fmt.Sscanf(parts[1], "%d", &id)
			if err != nil {
				fmt.Println("Invalid task ID")
				continue
			}
			tm.CompleteTask(id)

		case "delete":
			// "delete" command requires a task ID
			if len(parts) < 2 {
				fmt.Println("Usage: delete <task_id>")
				continue
			}
			var id int
			_, err := fmt.Sscanf(parts[1], "%d", &id)
			if err != nil {
				fmt.Println("Invalid task ID")
				continue
			}
			tm.DeleteTask(id)

		case "search":
			// "search" command requires a keyword
			if len(parts) < 2 {
				fmt.Println("Usage: search <keyword>")
				continue
			}
			keyword := strings.Join(parts[1:], " ")
			tm.SearchTasks(keyword)

		case "quit", "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command. Try: add, list, complete, delete, search, or quit")
		}
	}
}

// main is the entry point of the program
// CLI Apps: uses command-line flags and arguments to control program behavior
func main() {
	// Define command-line flags
	// These allow the program to be controlled from the command line
	addTaskFlag := flag.String("add", "", "Add a new task")
	listFlag := flag.Bool("list", false, "List all tasks")
	completeFlag := flag.Int("complete", 0, "Mark a task as completed (by ID)")
	deleteFlag := flag.Int("delete", 0, "Delete a task (by ID)")
	searchFlag := flag.String("search", "", "Search tasks by keyword")

	// Parse command-line arguments
	flag.Parse()

	// Create a TaskManager instance
	// Uses tasks.json in current directory to store tasks persistently
	manager := NewTaskManager("tasks.json")

	// Check which flag was provided and take appropriate action
	// Only one flag should be true/non-zero at a time

	switch {
	case *addTaskFlag != "":
		// User provided the -add flag
		manager.AddTask(*addTaskFlag)

	case *listFlag:
		// User provided the -list flag
		manager.ListTasks()

	case *completeFlag > 0:
		// User provided the -complete flag with a task ID
		manager.CompleteTask(*completeFlag)

	case *deleteFlag > 0:
		// User provided the -delete flag with a task ID
		manager.DeleteTask(*deleteFlag)

	case *searchFlag != "":
		// User provided the -search flag with a keyword
		manager.SearchTasks(*searchFlag)

	default:
		// No flags provided - start interactive mode
		manager.InteractiveMode()
	}
}
