module todoapp/webserver

go 1.24.2

require todoapp/task v0.0.0

replace todoapp/task => ../task

replace todoapp/middleware => ../middleware
