module todoapp/handlers

go 1.24.2

require (
	todoapp/middleware v0.0.0
	todoapp/task v0.0.0
)

require github.com/google/uuid v1.6.0 // indirect

replace todoapp/task => ../task

replace todoapp/middleware => ../middleware
