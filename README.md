# todoapp
A simple task management application that can be run as a CLI or an API.

-----

# Running the CLI
The CLI allows you to interact with the task manager directly from the command line.

## Steps to Run the CLI
1. Navigate to the cli directory:
``` shell
cd cli
```

2. Run the CLI:
``` shell
go run main.go
```

3. Use the CLI commands:
- *help*: Displays the list of available commands.
- *create*: Create a new task. Example:
``` shell
create -title Task1 -description Description1 -status NotStarted
```
- *list*: List all tasks.
```
- *exit*: Exist the CLI.

-----

# Running the API
The API provides HTTP endpoints to interact with the task manager.

## Steps to Run the API
1. Navigate to the root directory:
2. Run the API:
``` shell
go run main.go
```
3. Access the API:
- *Create Task*: <code>POST /create</code>
- *Get Tasks*: <code>GET /get</code>
- *Update Task*: <code>PUT /update</code>
- *Delete Task*: <code>DELETE /delete/{id}</code>
4. Use a tool like <code>curl</code> or <code>Postman</code> to interact with the API.

-----

# Example Usage
## CLI Example
``` bash
> go run main.go
CLI started. Type 'help' for commands.
> help
Available commands:
  get                                   - Retrieve and display all tasks
  create -title <title> -desc <description> -status <status> - Create a new task with the given details.
      Example: create -title Golang -desc Task1 -status NotStarted
  exit                                  - Exit the CLI
  help                                  - Show this help message
> create -title Task1 -description Description1 -status NotStarted
Task created successfully.
> list
1. Task 1 - Description 1 [NotStarted]
> exit
Exiting CLI...
```

## API example
``` bash
# Start the API
go run main.go

# Create a task
curl -X POST -H "Content-Type: application/json" -d '{"title":"Task 1","description":"Description 1","status":"NotStarted"}' http://localhost:8080/create

# Get all tasks
curl http://localhost:8080/get

# Update a task
curl -X PUT -H "Content-Type: application/json" -d '{"id":1,"title":"Updated Task 1","description":"Updated Description","status":"Completed"}' http://localhost:8080/update

# Delete a task
curl -X DELETE http://localhost:8080/delete/1
```

-----

# Dependencies
- Go 1.22 or higher
- Modules:
    - <code>todoapp/task</code>
    - <code>todoapp/files</code>
    - <code>todoapp/handlers</code>
    - <code>todoapp/middleware</code>
    - <code>todoapp/webserver</code>

-----

# Project Structure
```
todoapp/
├── cli/          # CLI implementation
├── task/         # Task management logic
├── files/        # File operations
├── handlers/     # HTTP handlers
├── middleware/   # Middleware for the API
├── webserver/    # Static and dynamic web pages
├── main.go       # API entry point
```