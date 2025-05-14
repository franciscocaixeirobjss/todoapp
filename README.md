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

4. Access the API:
- *Create Task*: <code>POST /create</code>
- *Get Tasks*: <code>GET /get</code>
- *Update Task*: <code>PUT /update</code>
- *Delete Task*: <code>DELETE /delete/{id}</code>

5. Use a tool like <code>curl</code> or <code>Postman</code> to interact with the API.

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
# Benchmarking Task

## Benchmark Actor Pattern
``` bash
> go test -benchmem -run=^$ -bench ^BenchmarkActorPattern$ todoapp/task

goos: windows
goarch: amd64
pkg: todoapp/task
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkActorPattern-8          1747880               678.8 ns/op           677 B/op          2 allocs/op
PASS
ok      todoapp/task    3.389s
```

## Benchmark Non Actor Pattern
``` bash
> go test -benchmem -run=^$ -bench ^BenchmarkNonActorPattern$ todoapp/task

goos: windows
goarch: amd64
pkg: todoapp/task
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkNonActorPattern-8       5565408               260.2 ns/op           585 B/op          1 allocs/op
PASS
ok      todoapp/task    2.017s
```

-----
# Benchmark Handlers

``` bash
>go test -benchmem -run=^ -bench ^Benchmark -count 3 todoapp/handlers

goos: windows
goarch: amd64
pkg: todoapp/handlers
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkUpdateHandlerNonActor-8                  286678              3514 ns/op            7077 B/op         33 allocs/op
BenchmarkUpdateHandlerNonActor-8                  313672              3673 ns/op            7078 B/op         33 allocs/op
BenchmarkUpdateHandlerNonActor-8                  317114              3612 ns/op            7078 B/op         33 allocs/op
BenchmarkUpdateHandlerActor_NoParallel-8          140314              8578 ns/op            6957 B/op         29 allocs/op
BenchmarkUpdateHandlerActor_NoParallel-8          136471              8578 ns/op            6958 B/op         29 allocs/op
BenchmarkUpdateHandlerActor_NoParallel-8          131121              8522 ns/op            6957 B/op         29 allocs/op
BenchmarkUpdateHandlerActor-8                     341443              3590 ns/op            7190 B/op         34 allocs/op
BenchmarkUpdateHandlerActor-8                     335078              3409 ns/op            7189 B/op         34 allocs/op
BenchmarkUpdateHandlerActor-8                     340693              3276 ns/op            7189 B/op         34 allocs/op
BenchmarkCreateHandlerNonActor-8                  460478              2184 ns/op            7641 B/op         33 allocs/op
BenchmarkCreateHandlerNonActor-8                  512782              2126 ns/op            7582 B/op         33 allocs/op
BenchmarkCreateHandlerNonActor-8                  558259              2191 ns/op            7660 B/op         33 allocs/op
BenchmarkCreateHandlerActor-8                     440749              2500 ns/op            7817 B/op         35 allocs/op
BenchmarkCreateHandlerActor-8                     693008              2596 ns/op            7782 B/op         35 allocs/op
BenchmarkCreateHandlerActor-8                     911541              2625 ns/op            7910 B/op         35 allocs/op
BenchmarkCreateHandlerActor_NoParallel-8          225532              5853 ns/op            6984 B/op         30 allocs/op
BenchmarkCreateHandlerActor_NoParallel-8          220810              7061 ns/op            8435 B/op         30 allocs/op
BenchmarkCreateHandlerActor_NoParallel-8          219613              5970 ns/op            6984 B/op         30 allocs/op
PASS
ok      todoapp/handlers        27.196s
```

-----

## Benchmark Unbuffered response channel and RequestChan Size 1M 
``` bash
>go test -benchmem -run=^ -bench ^Benchmark -count 3 todoapp/handlers

goos: windows
goarch: amd64
pkg: todoapp/handlers
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkUpdateHandlerActor-8             461100              3185 ns/op            7478 B/op         34 allocs/op
BenchmarkUpdateHandlerActor-8             430144              2698 ns/op            7500 B/op         34 allocs/op
BenchmarkUpdateHandlerActor-8             421582              2372 ns/op            7506 B/op         34 allocs/op
BenchmarkCreateHandlerActor-8             529372              2432 ns/op            8053 B/op         34 allocs/op
BenchmarkCreateHandlerActor-8             596853              3172 ns/op            7898 B/op         34 allocs/op
BenchmarkCreateHandlerActor-8             469543              2871 ns/op            8011 B/op         34 allocs/op
PASS
ok      todoapp/handlers        19.762s
```

-----

## Benchmark Buffered response channel and RequestChan Size 1M 
``` bash
>go test -benchmem -run=^ -bench ^Benchmark -count 3 todoapp/handlers

goos: windows
goarch: amd64
pkg: todoapp/handlers
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkUpdateHandlerActor-8             436520              2472 ns/op            7543 B/op         35 allocs/op
BenchmarkUpdateHandlerActor-8             334546              3623 ns/op            7638 B/op         35 allocs/op
BenchmarkUpdateHandlerActor-8             373936              3668 ns/op            7595 B/op         35 allocs/op
BenchmarkCreateHandlerActor-8             383052              3575 ns/op            8120 B/op         35 allocs/op
BenchmarkCreateHandlerActor-8             394144              3561 ns/op            8048 B/op         35 allocs/op
BenchmarkCreateHandlerActor-8             461997              2999 ns/op            8317 B/op         35 allocs/op
PASS
ok      todoapp/handlers        19.287s
```

-----

## Benchmark Unbuffered response channel and RequestChan Size 100 
``` bash
>go test -benchmem -run=^ -bench ^Benchmark -count 3 todoapp/handlers

goos: windows
goarch: amd64
pkg: todoapp/handlers
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkUpdateHandlerActor-8             200449              7218 ns/op            7192 B/op         34 allocs/op
BenchmarkUpdateHandlerActor-8             232488              6680 ns/op            7192 B/op         34 allocs/op
BenchmarkUpdateHandlerActor-8             214821              5908 ns/op            7192 B/op         34 allocs/op
BenchmarkCreateHandlerActor-8             388557              3313 ns/op            7813 B/op         34 allocs/op
BenchmarkCreateHandlerActor-8             947866              2570 ns/op            7814 B/op         34 allocs/op
BenchmarkCreateHandlerActor-8             913533              2886 ns/op            7807 B/op         34 allocs/op
PASS
ok      todoapp/handlers        16.279s
```

-----

## Benchmark Buffered response channel and RequestChan Size 100 
``` bash
>go test -benchmem -run=^ -bench ^Benchmark -count 3 todoapp/handlers

goos: windows
goarch: amd64
pkg: todoapp/handlers
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkUpdateHandlerActor-8             267734              4562 ns/op            7241 B/op         35 allocs/op
BenchmarkUpdateHandlerActor-8             265640              5040 ns/op            7241 B/op         35 allocs/op
BenchmarkUpdateHandlerActor-8             250174              5002 ns/op            7241 B/op         35 allocs/op
BenchmarkCreateHandlerActor-8             473222              2537 ns/op            7820 B/op         35 allocs/op
BenchmarkCreateHandlerActor-8             934215              2496 ns/op            7871 B/op         35 allocs/op
BenchmarkCreateHandlerActor-8             920209              2450 ns/op            7850 B/op         35 allocs/op
PASS
ok      todoapp/handlers        12.007s
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