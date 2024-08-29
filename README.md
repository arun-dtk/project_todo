# Todo application #
This is an application for CRUD operations on todo items. This is built using Golang and Gin framework.

## Project Setup / Installation
Clone the project code.  
Run `cd project_todo` to move to project folder in terminal.  
Install all dependencies.  
Set up a Postgres database and keep the connection url in the env file on your project home directory.
Use REST Client extension in VS code for testing APIs.

## Running the project
Run `go run .` to run the project and access it via browser on `localhost:8080`

## Running the API benchmark
Open a new terminal and go to benchmark folder using `cd benchmark` & run using `go run benchmark.go` 

## Project creation steps
1. `go mod init project_todo` This will create go.mod file.
2. `go get -u github.com/gin-gonic/gin` This will add gin framework dependency to the project.

## env file items expected. Please fill in proper details to run the project

DB_HOST=""

DB_PORT=5432

DB_USER=""

DB_PASSWORD=""

DB_NAME=""

JWT_SECRET=""

