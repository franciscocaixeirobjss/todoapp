module todoapp/backend

go 1.24.2

require (
	todoapp/files v0.0.0
	todoapp/handlers v0.0.0
	todoapp/logging v0.0.0
	todoapp/middleware v0.0.0
	todoapp/task v0.0.0
	todoapp/webserver v0.0.0
)

require github.com/google/uuid v1.6.0 // indirect

replace todoapp/task => ../task

replace todoapp/files => ../files

replace todoapp/handlers => ../handlers

replace todoapp/webserver => ../webserver

replace todoapp/logging => ../logging

replace todoapp/middleware => ../middleware
