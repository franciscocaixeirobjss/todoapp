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

## Benchmark Create 
``` bash
> go test -benchmem -run=^ -bench ^BenchmarkCreate -count 3  todoapp/task

goos: windows
goarch: amd64
pkg: todoapp/task
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkCreateActorPattern-8            2030899               524.4 ns/op           601 B/op          2 allocs/op
BenchmarkCreateActorPattern-8            1579746               666.5 ns/op           532 B/op          2 allocs/op
BenchmarkCreateActorPattern-8            2094079               618.8 ns/op           719 B/op          2 allocs/op
BenchmarkCreateNonActorPattern-8         5433957               198.7 ns/op           599 B/op          1 allocs/op
BenchmarkCreateNonActorPattern-8         8821278               221.3 ns/op           577 B/op          1 allocs/op
BenchmarkCreateNonActorPattern-8         8959837               212.5 ns/op           569 B/op          1 allocs/op
PASS
ok      todoapp/task    16.088s
```

## Benchmark Update
``` bash
> go test -benchmem -run=^ -bench ^BenchmarkUpdate -count 3  todoapp/task

goos: windows
goarch: amd64
pkg: todoapp/task
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkUpdateActorPattern-8            2729719               439.6 ns/op           136 B/op          2 allocs/op
BenchmarkUpdateActorPattern-8            2590282               496.7 ns/op           136 B/op          2 allocs/op
BenchmarkUpdateActorPattern-8            2438349               440.4 ns/op           136 B/op          2 allocs/op
BenchmarkUpdateNonActorPattern-8        12953227                91.17 ns/op           24 B/op          1 allocs/op
BenchmarkUpdateNonActorPattern-8        13277992                87.23 ns/op           24 B/op          1 allocs/op
BenchmarkUpdateNonActorPattern-8        13183190                90.77 ns/op           24 B/op          1 allocs/op
PASS
ok      todoapp/task    11.054s
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
# Unit tests

## Task
``` bash
> go test -race -v

=== RUN   TestCreateTask
--- PASS: TestCreateTask (0.00s)
=== RUN   TestGetTasks
--- PASS: TestGetTasks (0.00s)
=== RUN   TestUpdateTask
--- PASS: TestUpdateTask (0.00s)
=== RUN   TestDeleteTask
--- PASS: TestDeleteTask (0.00s)
=== RUN   TestConvertStringToStatusID
=== RUN   TestConvertStringToStatusID/NotStarted
=== RUN   TestConvertStringToStatusID/_Not_Started_
=== RUN   TestConvertStringToStatusID/Started
=== RUN   TestConvertStringToStatusID/Completed
=== RUN   TestConvertStringToStatusID/Invalid_Status
=== RUN   TestConvertStringToStatusID/InvalidStatus
--- PASS: TestConvertStringToStatusID (0.00s)
    --- PASS: TestConvertStringToStatusID/NotStarted (0.00s)
    --- PASS: TestConvertStringToStatusID/_Not_Started_ (0.00s)
    --- PASS: TestConvertStringToStatusID/Started (0.00s)
    --- PASS: TestConvertStringToStatusID/Completed (0.00s)
    --- PASS: TestConvertStringToStatusID/Invalid_Status (0.00s)
    --- PASS: TestConvertStringToStatusID/InvalidStatus (0.00s)
PASS
ok      todoapp/task    2.988s
```

## Handlers
``` bash
> go test -race -v

=== RUN   TestCreateHandler_ServiceUnavailable
--- PASS: TestCreateHandler_ServiceUnavailable (0.00s)
=== RUN   TestCreateHandler_Parallel
=== RUN   TestCreateHandler_Parallel/Task_0
=== PAUSE TestCreateHandler_Parallel/Task_0
...
=== CONT  TestCreateHandler_Parallel/Task_0
=== CONT  TestCreateHandler_Parallel/Task_3
--- PASS: TestCreateHandler_Parallel (0.02s)
    --- PASS: TestCreateHandler_Parallel/Task_1 (0.00s)
    --- PASS: TestCreateHandler_Parallel/Task_53 (0.00s)
    ...
    --- PASS: TestCreateHandler_Parallel/Task_0 (0.00s)
    --- PASS: TestCreateHandler_Parallel/Task_3 (0.00s)
=== RUN   TestCreateHandler_Goroutine_Parallel
--- PASS: TestCreateHandler_Goroutine_Parallel (0.00s)
=== RUN   TestCreateHandler
=== RUN   TestCreateHandler/valid_create_request
=== RUN   TestCreateHandler/status_bad_request_-_invalid_json_format
=== RUN   TestCreateHandler/method_not_allowed_-_get_request_instead_of_post
--- PASS: TestCreateHandler (0.00s)
    --- PASS: TestCreateHandler/valid_create_request (0.00s)
    --- PASS: TestCreateHandler/status_bad_request_-_invalid_json_format (0.00s)
    --- PASS: TestCreateHandler/method_not_allowed_-_get_request_instead_of_post (0.00s)
=== RUN   TestGetHandler
=== RUN   TestGetHandler/valid_get_request
=== RUN   TestGetHandler/method_not_allowed_-_post_request_instead_of_get
--- PASS: TestGetHandler (0.00s)
    --- PASS: TestGetHandler/valid_get_request (0.00s)
    --- PASS: TestGetHandler/method_not_allowed_-_post_request_instead_of_get (0.00s)
=== RUN   TestTaskActor_Concurrency
--- PASS: TestTaskActor_Concurrency (0.00s)
=== RUN   TestTaskActor_ConcurrentUpdate
--- PASS: TestTaskActor_ConcurrentUpdate (0.00s)
=== RUN   TestUpdateHandler
=== RUN   TestUpdateHandler/valid_update_request
=== RUN   TestUpdateHandler/not_found_-_non-existing_id
=== RUN   TestUpdateHandler/method_not_allowed_-_post_request_instead_of_put
--- PASS: TestUpdateHandler (0.00s)
    --- PASS: TestUpdateHandler/valid_update_request (0.00s)
    --- PASS: TestUpdateHandler/not_found_-_non-existing_id (0.00s)
    --- PASS: TestUpdateHandler/method_not_allowed_-_post_request_instead_of_put (0.00s)
=== RUN   TestDeleteHandler
=== RUN   TestDeleteHandler/valid_delete_request
=== RUN   TestDeleteHandler/not_found_-_non-existing_id
=== RUN   TestDeleteHandler/not_found_-_already_deleted_task
=== RUN   TestDeleteHandler/bad_request_-_invalid_task_ID
=== RUN   TestDeleteHandler/method_not_allowed_-_post_request_instead_of_delete
--- PASS: TestDeleteHandler (0.00s)
    --- PASS: TestDeleteHandler/valid_delete_request (0.00s)
    --- PASS: TestDeleteHandler/not_found_-_non-existing_id (0.00s)
    --- PASS: TestDeleteHandler/not_found_-_already_deleted_task (0.00s)
    --- PASS: TestDeleteHandler/bad_request_-_invalid_task_ID (0.00s)
    --- PASS: TestDeleteHandler/method_not_allowed_-_post_request_instead_of_delete (0.00s)
PASS
ok      todoapp/handlers        4.256s
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